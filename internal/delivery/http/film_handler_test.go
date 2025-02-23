package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	filmHttp "go-films-api/internal/delivery/http"
	"go-films-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFilmService struct {
	mock.Mock
}

func (m *MockFilmService) ListFilms(title, genre string, releaseDate time.Time) ([]domain.Film, error) {
	args := m.Called(title, genre, releaseDate)
	if films, ok := args.Get(0).([]domain.Film); ok {
		return films, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFilmService) GetFilmDetails(id uint) (*domain.Film, error) {
	args := m.Called(id)
	if film, ok := args.Get(0).(*domain.Film); ok {
		return film, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetFilms_NoFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockFilmService)
	filmHandler := filmHttp.NewFilmHandler(mockService)

	r := gin.Default()
	r.GET("/films", filmHandler.GetFilms)

	expectedFilms := []domain.Film{
		{ID: 1, Title: "Film One", Genre: "Action"},
		{ID: 2, Title: "Film Two", Genre: "Drama"},
	}

	mockService.On("ListFilms", "", "", time.Time{}).Return(expectedFilms, nil)

	req, _ := http.NewRequest("GET", "/films", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "should return 200")

	assert.Contains(t, w.Body.String(), "Film One")
	assert.Contains(t, w.Body.String(), "Film Two")

	mockService.AssertExpectations(t)
}

func TestGetFilms_WithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockFilmService)
	filmHandler := filmHttp.NewFilmHandler(mockService)

	r := gin.Default()
	r.GET("/films", filmHandler.GetFilms)

	date, _ := time.Parse("2006-01-02", "2023-01-01")
	expectedFilms := []domain.Film{
		{ID: 10, Title: "Action Film", Genre: "Action"},
	}

	mockService.On("ListFilms", "Action", "Action", date).
		Return(expectedFilms, nil)

	req, _ := http.NewRequest("GET", "/films?title=Action&genre=Action&release_date=2023-01-01", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "should return 200")

	assert.Contains(t, w.Body.String(), "Action Film")
	mockService.AssertExpectations(t)
}

func TestGetFilms_InvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockFilmService)
	filmHandler := filmHttp.NewFilmHandler(mockService)

	r := gin.Default()
	r.GET("/films", filmHandler.GetFilms)

	req, _ := http.NewRequest("GET", "/films?release_date=invalid-date", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "should return 400 for bad date")
	assert.Contains(t, w.Body.String(), "invalid release_date format")
}

func TestGetFilms_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockFilmService)
	filmHandler := filmHttp.NewFilmHandler(mockService)

	r := gin.Default()
	r.GET("/films", filmHandler.GetFilms)

	mockService.On("ListFilms", "", "", time.Time{}).
		Return(nil, fmt.Errorf("some db error"))

	req, _ := http.NewRequest("GET", "/films", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to fetch films")

	mockService.AssertExpectations(t)
}

func TestGetFilmDetails_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockFilmService)
	filmHandler := filmHttp.NewFilmHandler(mockService)

	r := gin.Default()
	r.GET("/films/:id", filmHandler.GetFilmDetails)

	expectedFilm := &domain.Film{
		ID:    1,
		Title: "My Film",
		User:  domain.User{ID: 2, Username: "creatoruser"},
	}

	mockService.
		On("GetFilmDetails", uint(1)).
		Return(expectedFilm, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/films/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My Film")
	assert.Contains(t, w.Body.String(), "creatoruser")

	mockService.AssertExpectations(t)
}

func TestGetFilmDetails_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockFilmService)
	filmHandler := filmHttp.NewFilmHandler(mockService)

	r := gin.Default()
	r.GET("/films/:id", filmHandler.GetFilmDetails)

	mockService.
		On("GetFilmDetails", uint(99)).
		Return(nil, errors.New("film not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/films/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "film not found")

	mockService.AssertExpectations(t)
}
