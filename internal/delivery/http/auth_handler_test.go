package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authHttp "go-films-api/internal/delivery/http"
	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
	"go-films-api/internal/usecase"
)

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(repository.MockUserRepository)
	userService := usecase.NewUserService(mockRepo)
	authHandler := authHttp.NewAuthHandler(userService)

	r := gin.Default()
	r.POST("/register", authHandler.Register)

	mockRepo.On("GetUserByUsername", "newuser").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	body := `{"username":"newuser","password":"secret"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "user registered successfully", resp["message"])
	mockRepo.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(repository.MockUserRepository)
	userService := usecase.NewUserService(mockRepo)
	authHandler := authHttp.NewAuthHandler(userService)

	r := gin.Default()
	r.POST("/login", authHandler.Login)

	// Provide a known hashed password matching "secret"
	hashed := "$2a$10$1fybhpdIC527ODopk5/FLu5L5o60g.2p1NGd7Zso75iv.R4siZm3e"
	user := &domain.User{ID: 1, Username: "alex", Password: hashed}

	mockRepo.On("GetUserByUsername", "alex").Return(user, nil)

	body := `{"username":"alex","password":"secret"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["token"], "Expected a token in response")
	mockRepo.AssertExpectations(t)
}
