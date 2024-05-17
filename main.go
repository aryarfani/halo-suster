package main

import (
	"context"
	"eniqilo-store/config/storage"
	"eniqilo-store/db"
	"eniqilo-store/router"
	"eniqilo-store/utils"
	"log"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	db.Connect()
	storage.New()

	utils.InitValidator()

	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      "Cats Social",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(logger.New())
	app.Use(healthcheck.New())

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))

	ctx := context.Background()
	wait := utils.GracefulShutdown(ctx, 5*time.Second, map[string]utils.Operation{
		"postgre": func(ctx context.Context) error {
			return db.Close()
		},
		"fiber-server": func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	<-wait
}
