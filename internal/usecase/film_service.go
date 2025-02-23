package usecase

import (
	"time"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
)

type FilmService interface {
	ListFilms(title, genre string, releaseDate time.Time) ([]domain.Film, error)
}

type filmService struct {
	filmRepo repository.FilmRepository
}

func NewFilmService(repo repository.FilmRepository) FilmService {
	return &filmService{filmRepo: repo}
}

func (s *filmService) ListFilms(title, genre string, releaseDate time.Time) ([]domain.Film, error) {
	filters := repository.FilmFilters{
		Title:       title,
		Genre:       genre,
		ReleaseDate: releaseDate,
	}
	return s.filmRepo.FindFilms(filters)
}
