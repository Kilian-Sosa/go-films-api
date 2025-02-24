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

type CreateFilmRequest struct {
	Title       string `json:"title" binding:"required"`
	Director    string `json:"director"`
	ReleaseDate string `json:"release_date"`
	Cast        string `json:"cast"`
	Genre       string `json:"genre"`
	Synopsis    string `json:"synopsis"`
}

type UpdateFilmRequest struct {
	Title       *string `json:"title"`
	Director    *string `json:"director"`
	ReleaseDate *string `json:"release_date"`
	Cast        *string `json:"cast"`
	Genre       *string `json:"genre"`
	Synopsis    *string `json:"synopsis"`
}

func NewFilmHandler(fs usecase.FilmService) *FilmHandler {
	return &FilmHandler{filmService: fs}
}

// GetFilms godoc
// @Summary Get a list of films
// @Description Retrieves a list of films, optionally filtered by title, genre, and release date.
// @Tags films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param title query string false "Film title"
// @Param genre query string false "Film genre"
// @Param release_date query string false "Film release date (YYYY-MM-DD)"
// @Success 200 {array} domain.Film
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /films [get]
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

// GetFilmDetails godoc
// @Summary Get details of a specific film
// @Description Retrieves the details of a film by ID, including the creator user.
// @Tags films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Film ID"
// @Success 200 {object} domain.Film
// @Failure 400 {object} map[string]string "Invalid Film ID"
// @Failure 404 {object} map[string]string "Film not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /films/{id} [get]
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

// CreateFilm godoc
// @Summary Create a new film
// @Description Adds a new film to the database, linked to the authenticated user.
// @Tags films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param film body CreateFilmRequest true "Film details"
// @Success 201 {object} domain.Film
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 409 {object} map[string]string "Film already exists"
// @Router /films [post]
func (h *FilmHandler) CreateFilm(c *gin.Context) {
	var req CreateFilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var rd time.Time
	var err error
	if req.ReleaseDate != "" {
		rd, err = time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release_date format, expected YYYY-MM-DD"})
			return
		}
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	film, createErr := h.filmService.CreateFilm(
		req.Title,
		req.Director,
		req.Cast,
		req.Genre,
		req.Synopsis,
		rd,
		userIDValue.(uint),
	)
	if createErr != nil {
		c.JSON(http.StatusConflict, gin.H{"error": createErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, film)
}

// UpdateFilm godoc
// @Summary Update a film
// @Description Updates the details of a film, only allowed for the creator user.
// @Tags films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Film ID"
// @Param film body UpdateFilmRequest true "Film details"
// @Success 200 {object} domain.Film
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 403 {object} map[string]string "Forbidden: only creator can update this film"
// @Failure 404 {object} map[string]string "Film not found"
// @Failure 409 {object} map[string]string "Could not update film"
// @Router /films/{id} [put]
func (h *FilmHandler) UpdateFilm(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}
	filmID := uint(id64)

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	var req UpdateFilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var releaseDatePtr *time.Time
	if req.ReleaseDate != nil && *req.ReleaseDate != "" {
		rd, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release_date format, expected YYYY-MM-DD"})
			return
		}
		releaseDatePtr = &rd
	}

	data := usecase.UpdateFilmData{
		Title:       req.Title,
		Director:    req.Director,
		Cast:        req.Cast,
		Genre:       req.Genre,
		Synopsis:    req.Synopsis,
		ReleaseDate: releaseDatePtr,
	}

	updated, err := h.filmService.UpdateFilm(filmID, userID, data)
	if err != nil {
		switch err.Error() {
		case "film not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "film not found"})
		case "forbidden: only creator can update this film":
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only creator can update this film"})
		default:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteFilm godoc
// @Summary Delete a film
// @Description Deletes a film from the database, only allowed for the creator user.
// @Tags films
// @Security BearerAuth
// @Param id path int true "Film ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid Film ID"
// @Failure 403 {object} map[string]string "Forbidden: only creator can delete this film"
// @Failure 404 {object} map[string]string "Film not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /films/{id} [delete]
func (h *FilmHandler) DeleteFilm(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}
	filmID := uint(id64)

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	err = h.filmService.DeleteFilm(filmID, userID)
	if err != nil {
		switch err.Error() {
		case "film not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "film not found"})
		case "forbidden: only creator can delete this film":
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only creator can delete this film"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete film"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
