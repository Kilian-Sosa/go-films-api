package http

import (
	"net/http"
	"strconv"
	"time"

	"go-films-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FilmHandler struct {
	filmService usecase.FilmService
}

func NewFilmHandler(fs usecase.FilmService) *FilmHandler {
	return &FilmHandler{filmService: fs}
}

func (h *FilmHandler) GetFilms(c *gin.Context) {
	title := c.Query("title")
	genre := c.Query("genre")

	releaseDateStr := c.Query("release_date")
	var releaseDate time.Time
	var err error
	if releaseDateStr != "" {
		releaseDate, err = time.Parse("2006-01-02", releaseDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release_date format, expected YYYY-MM-DD"})
			return
		}
	}

	films, err := h.filmService.ListFilms(title, genre, releaseDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch films"})
		return
	}

	c.JSON(http.StatusOK, films)
}

func (h *FilmHandler) GetFilmDetails(c *gin.Context) {
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}
	filmID := uint(id64)

	film, err := h.filmService.GetFilmDetails(filmID)
	if err != nil {
		if err.Error() == "film not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "film not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve film details"})
		}
		return
	}

	c.JSON(http.StatusOK, film)
}
