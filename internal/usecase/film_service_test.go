package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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

func TestGetFilmDetails_Found(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	service := usecase.NewFilmService(mockRepo)

	expectedFilm := &domain.Film{
		ID:    1,
		Title: "My Film",
		User:  domain.User{ID: 2, Username: "creator"},
	}

	mockRepo.On("GetFilmByID", uint(1)).Return(expectedFilm, nil)

	film, err := service.GetFilmDetails(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), film.ID)
	assert.Equal(t, "creator", film.User.Username)
	mockRepo.AssertExpectations(t)
}

func TestGetFilmDetails_NotFound(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	service := usecase.NewFilmService(mockRepo)

	mockRepo.On("GetFilmByID", uint(99)).Return(nil, nil)

	film, err := service.GetFilmDetails(99)
	assert.Nil(t, film)
	assert.EqualError(t, err, "film not found")
	mockRepo.AssertExpectations(t)
}

func TestCreateFilm_Success(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	filmService := usecase.NewFilmService(mockRepo)

	mockRepo.On("CreateFilm", mock.AnythingOfType("*domain.Film")).
		Return(nil).
		Run(func(args mock.Arguments) {
			// Simulate setting an auto-increment ID
			arg := args.Get(0).(*domain.Film)
			arg.ID = 100
		})

	res, err := filmService.CreateFilm(
		"Unique Title", "Director", "Cast", "Action", "Some synopsis", time.Time{}, 1,
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(100), res.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateFilm_DuplicateTitle(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	filmService := usecase.NewFilmService(mockRepo)

	mockRepo.On("CreateFilm", mock.Anything).
		Return(fmt.Errorf("film with title 'Duplicate' already exists"))

	res, err := filmService.CreateFilm("Duplicate", "", "", "", "", time.Time{}, 1)
	assert.Nil(t, res)
	assert.EqualError(t, err, "film with title 'Duplicate' already exists")
	mockRepo.AssertExpectations(t)
}

func TestCreateFilm_EmptyTitle(t *testing.T) {
	mockRepo := new(repository.MockFilmRepository)
	filmService := usecase.NewFilmService(mockRepo)

	res, err := filmService.CreateFilm("", "Dir", "Cast", "Genre", "Synopsis", time.Time{}, 1)
	assert.Nil(t, res)
	assert.EqualError(t, err, "title is required")
	mockRepo.AssertNotCalled(t, "CreateFilm", mock.Anything)
}
