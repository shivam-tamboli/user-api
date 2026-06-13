// @title           User API
// @version         1.0
// @description     A RESTful API to manage users with dynamic age calculation.
// @BasePath        /

package main

import (
	"context"
	"log"
	"user-api/config"
	_ "user-api/docs"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer conn.Close(ctx)

	logger.Log.Info("connected to database")

	if err := repository.RunMigration(ctx, conn); err != nil {
		logger.Log.Fatal("migration failed", zap.Error(err))
	}
	logger.Log.Info("migrations applied")

	queries := repository.NewRepository(conn)
	userService := service.NewUserService(queries)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Setup(app, userHandler)

	logger.Log.Info("server starting", zap.String("port", cfg.Port))
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Log.Fatal("server failed", zap.Error(err))
	}
}
