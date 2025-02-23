package usecase

import (
	"errors"
	"time"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
)

type FilmService interface {
	ListFilms(title, genre string, releaseDate time.Time) ([]domain.Film, error)
	GetFilmDetails(id uint) (*domain.Film, error)
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

func (s *filmService) GetFilmDetails(id uint) (*domain.Film, error) {
	film, err := s.filmRepo.GetFilmByID(id)
	if err != nil {
		return nil, err
	}
	if film == nil {
		return nil, errors.New("film not found")
	}
	return film, nil
}
