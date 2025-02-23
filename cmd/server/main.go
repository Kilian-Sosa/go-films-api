package main

import (
	"fmt"
	"go-films-api/internal/delivery/http"
	"go-films-api/internal/delivery/http/middleware"
	"go-films-api/internal/repository"
	"go-films-api/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	userRepo := repository.NewUserRepositoryGorm(db)
	userService := usecase.NewUserService(userRepo)

	authHandler := http.NewAuthHandler(userService)

	filmRepo := repository.NewFilmRepositoryGorm(db)
	filmService := usecase.NewFilmService(filmRepo)
	filmHandler := http.NewFilmHandler(filmService)

	authMiddleware := middleware.JWTMiddleware()

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("/films", filmHandler.GetFilms)
		protected.GET("/films/:id", filmHandler.GetFilmDetails)
		protected.POST("/films", filmHandler.CreateFilm)
		protected.PUT("/films/:id", filmHandler.UpdateFilm)
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
