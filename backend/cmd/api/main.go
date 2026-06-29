package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"taladsync/internal/auth"
	"taladsync/internal/config"
	"taladsync/internal/database"
	"taladsync/internal/handler"
	"taladsync/internal/middleware"
	"taladsync/internal/repository"
	"taladsync/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err := database.SeedIfEmpty(db); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenSvc := auth.NewTokenService([]byte(cfg.JWTSecret))
	authUC := usecase.NewAuthUsecase(userRepo, tokenSvc)
	authHandler := handler.NewAuthHandler(authUC)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	api.Post("/auth/mock-login", authHandler.MockLogin)

	protected := api.Group("/auth", middleware.Auth(tokenSvc))
	protected.Get("/me", authHandler.Me)

	addr := ":" + cfg.AppPort
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
