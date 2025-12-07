package main

import (
	"log"
	adapterTeacher "myapp/adapter/teacher"
	coreTeacher "myapp/core/teacher"
	"myapp/middlware"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	repo := adapterTeacher.NewGormTeacherRepo(db)
	service := coreTeacher.NewTeacherApp(coreTeacher.WithTeacherRepo(repo))
	handler := adapterTeacher.NewHttpTeacherHandler(service)

	app := fiber.New()
	app.Use(middlware.LoggingMiddleware)
	app.Get("/teachers/:teacher_id/students", handler.GetAllStudents)
	app.Post("/teachers/login", handler.Login)
	app.Post("/teachers/register", handler.Register)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}
	app.Listen(":" + port)
}
