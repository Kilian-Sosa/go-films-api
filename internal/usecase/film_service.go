package usecase

import (
	"errors"
	"fmt"
	"time"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
)

type FilmService interface {
	ListFilms(title, genre string, releaseDate time.Time) ([]domain.Film, error)
	GetFilmDetails(id uint) (*domain.Film, error)
	CreateFilm(title, director, cast, genre, synopsis string, releaseDate time.Time, userID uint) (*domain.Film, error)
	UpdateFilm(id, userID uint, data UpdateFilmData) (*domain.Film, error)
}

type UpdateFilmData struct {
	Title       *string
	Director    *string
	ReleaseDate *time.Time
	Cast        *string
	Genre       *string
	Synopsis    *string
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

func (s *filmService) CreateFilm(
	title, director, cast, genre, synopsis string,
	releaseDate time.Time,
	userID uint,
) (*domain.Film, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	film := &domain.Film{
		UserID:      userID,
		Title:       title,
		Director:    director,
		ReleaseDate: releaseDate,
		Cast:        cast,
		Genre:       genre,
		Synopsis:    synopsis,
	}

	if err := s.filmRepo.CreateFilm(film); err != nil {
		return nil, err
	}

	return film, nil
}

func (s *filmService) UpdateFilm(id, userID uint, data UpdateFilmData) (*domain.Film, error) {
	film, err := s.filmRepo.GetFilmByID(id)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if film == nil {
		return nil, errors.New("film not found")
	}

	if film.UserID != userID {
		return nil, errors.New("forbidden: only creator can update this film")
	}

	if data.Title != nil {
		film.Title = *data.Title
	}
	if data.Director != nil {
		film.Director = *data.Director
	}
	if data.ReleaseDate != nil {
		film.ReleaseDate = *data.ReleaseDate
	}
	if data.Cast != nil {
		film.Cast = *data.Cast
	}
	if data.Genre != nil {
		film.Genre = *data.Genre
	}
	if data.Synopsis != nil {
		film.Synopsis = *data.Synopsis
	}

	if err := s.filmRepo.UpdateFilm(film); err != nil {
		return nil, err
	}

	return film, nil
}
