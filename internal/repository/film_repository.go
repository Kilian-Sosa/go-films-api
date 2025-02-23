package repository

import (
	"time"

	"go-films-api/internal/domain"

	"gorm.io/gorm"
)

type FilmRepository interface {
	FindFilms(filters FilmFilters) ([]domain.Film, error)
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
