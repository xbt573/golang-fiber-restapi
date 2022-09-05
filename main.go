package main

import (
	"fmt"
	"restapi/database"
	"restapi/handlers"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	var args struct {
		Port             int    `arg:"-p,--port,env:PORT" help:"Port to start server on" default:"3000"`
		Prod             bool   `arg:"--prod,env:PROD" help:"Production mode" default:"false"`
		ConnectionString string `arg:"required,--db,env:DATABASE_URL" help:"Database connection string"`
		DatabaseName     string `arg:"--dbname,env:DATABASE_NAME" help:"Database name (optional)" default:"playground"`
	}
	arg.MustParse(&args)

	err := database.Connect(args.ConnectionString, args.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Disconnect()

	app := fiber.New(fiber.Config{
		Prefork: args.Prod,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/tasks", handlers.GetTasks)
	v1.Get("/tasks/:id", handlers.GetTask)

	v1.Post("/tasks", handlers.CreateTask)

	v1.Patch("/tasks/:id", handlers.UpdateTask)

	v1.Delete("/tasks/:id", handlers.DeleteTask)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", args.Port)))
}
