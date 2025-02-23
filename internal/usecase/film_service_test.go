package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
	"go-films-api/internal/usecase"
)

func TestListFilms_NoFilters(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	filmService := usecase.NewFilmService(mockRepo)

	expectedFilms := []domain.Film{
		{ID: 1, Title: "Film One", Genre: "Action"},
		{ID: 2, Title: "Film Two", Genre: "Drama"},
	}

	mockRepo.On("FindFilms", repository.FilmFilters{}).
		Return(expectedFilms, nil)

	films, err := filmService.ListFilms("", "", time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(films))
	assert.Equal(t, "Film One", films[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestListFilms_WithTitleFilter(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	filmService := usecase.NewFilmService(mockRepo)

	expectedFilms := []domain.Film{
		{ID: 3, Title: "Matrix Reloaded", Genre: "Sci-Fi"},
	}

	filters := repository.FilmFilters{
		Title: "Matrix",
	}

	mockRepo.On("FindFilms", filters).
		Return(expectedFilms, nil)

	films, err := filmService.ListFilms("Matrix", "", time.Time{})
	assert.NoError(t, err)
	assert.Len(t, films, 1)
	assert.Equal(t, "Matrix Reloaded", films[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestListFilms_WithGenreAndDate(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	filmService := usecase.NewFilmService(mockRepo)

	date, _ := time.Parse("2006-01-02", "2023-01-01")
	filters := repository.FilmFilters{
		Genre:       "Action",
		ReleaseDate: date,
	}

	expectedFilms := []domain.Film{
		{ID: 4, Title: "Action Film 2023", Genre: "Action"},
	}
	mockRepo.On("FindFilms", filters).
		Return(expectedFilms, nil)

	films, err := filmService.ListFilms("", "Action", date)
	assert.NoError(t, err)
	assert.Len(t, films, 1)
	assert.Equal(t, uint(4), films[0].ID)
	mockRepo.AssertExpectations(t)
}
