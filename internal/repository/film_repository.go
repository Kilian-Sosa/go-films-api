package repository

import (
	"errors"
	"fmt"
	"time"

	"go-films-api/internal/domain"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type FilmRepository interface {
	FindFilms(filters FilmFilters) ([]domain.Film, error)
	GetFilmByID(id uint) (*domain.Film, error)
	CreateFilm(film *domain.Film) error
}

type filmRepositoryGorm struct {
	db *gorm.DB
}

func NewFilmRepositoryGorm(db *gorm.DB) FilmRepository {
	return &filmRepositoryGorm{db: db}
}

type FilmFilters struct {
	Title       string
	Genre       string
	ReleaseDate time.Time
}

func (r *filmRepositoryGorm) FindFilms(filters FilmFilters) ([]domain.Film, error) {
	query := r.db.Model(&domain.Film{})

	if filters.Title != "" {
		query = query.Where("title LIKE ?", "%"+filters.Title+"%")
	}
	if filters.Genre != "" {
		query = query.Where("genre = ?", filters.Genre)
	}

	if !filters.ReleaseDate.IsZero() {
		query = query.Where("release_date = ?", filters.ReleaseDate)
	}

	var films []domain.Film
	if err := query.Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}

func (r *filmRepositoryGorm) GetFilmByID(id uint) (*domain.Film, error) {
	var film domain.Film
	err := r.db.Preload("User").First(&film, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not get film: %w", err)
	}
	return &film, nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlError *mysql.MySQLError
	return errors.As(err, &mysqlError) && mysqlError.Number == 1062
}

func (r *filmRepositoryGorm) CreateFilm(film *domain.Film) error {
	if err := r.db.Create(film).Error; err != nil {
		// Check if it's a duplicate key error on Title
		if isDuplicateKeyError(err) {
			return fmt.Errorf("film with title '%s' already exists", film.Title)
		}
		return err
	}
	return nil
}
