package usecase

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
)

type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, time.Time, error)
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

// Validation constants
const (
	PasswordMinLen = 6
	PasswordMaxLen = 20
)

// Regex for username: start with letter, then alphanumeric
var usernameRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

func (s *userService) Register(username, password string) error {
	if !usernameRegex.MatchString(username) {
		return errors.New("username must start with a letter and contain only alphanumeric characters")
	}

	if len(password) < PasswordMinLen || len(password) > PasswordMaxLen {
		return fmt.Errorf("password must be between %d and %d characters", PasswordMinLen, PasswordMaxLen)
	}

	existing, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	if existing != nil {
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
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err
	}

	return nil
}

func (s *userService) Login(username, password string) (string, time.Time, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("repository error: %w", err)
	}

	if user == nil {
		return "", time.Time{}, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", time.Time{}, errors.New("invalid username or password")
	}

	expirationTime := time.Now().Add(time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expirationTime.Unix(),
	})

	signedToken, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("could not sign token: %w", err)
	}

	return signedToken, expirationTime, nil
}
