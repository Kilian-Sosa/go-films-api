package usecase

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
)

type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
		jwtKey:   []byte(os.Getenv("JWT_SECRET")),
	}
}

func (s *userService) Register(username, password string) error {
	existingUser, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	if existingUser != nil {
		return errors.New("username already taken")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	newUser := &domain.User{
		Username: username,
		Password: string(hashedPass),
	}

	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("repository error: %w", err)
	}
	if user == nil {
		return "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signedToken, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return signedToken, nil
}
