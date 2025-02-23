package repository

import (
	"github.com/stretchr/testify/mock"

	"go-films-api/internal/domain"
)

type MockFilmRepository struct {
	mock.Mock
}

func (m *MockFilmRepository) FindFilms(filters FilmFilters) ([]domain.Film, error) {
	args := m.Called(filters)
	if films, ok := args.Get(0).([]domain.Film); ok {
		return films, args.Error(1)
	}
	return nil, args.Error(1)
}
