package main

import (
	_ "embed"
	"log"
	adapterClass "myapp/adapter/class"
	adapterTeacher "myapp/adapter/teacher"
	coreClass "myapp/core/class"
	coreTeacher "myapp/core/teacher"
	"myapp/middleware"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed databasetable.sql
var schemaSQL string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	if err := migrateDatabase(db); err != nil {
		log.Fatal(err)
	}

	teacherRepo := adapterTeacher.NewGormTeacherRepo(db)
	teacherService := coreTeacher.NewTeacherApp(
		coreTeacher.WithTeacherRepo(teacherRepo),
		coreTeacher.WithJWTSecret(os.Getenv("JWT_SECRET")),
	)
	teacherHandler := adapterTeacher.NewHttpTeacherHandler(teacherService)

	classRepo := adapterClass.NewGormClassRepo(db)
	classService := coreClass.NewClassApp(
		coreClass.WithClassRepo(classRepo),
	)
	classHandler := adapterClass.NewHttpClassHandler(classService)

	app := fiber.New()
	app.Use(middleware.LoggingMiddleware)
	app.Get("/teachers/:teacher_id/students", teacherHandler.GetAllStudents)
	app.Post("/teachers/login", teacherHandler.Login)
	app.Post("/teachers/register", teacherHandler.Register)
	app.Get("/teachers/:teacher_id/classes", classHandler.GetAllClasses)
	app.Post("/teachers/:teacher_id/classes", classHandler.CreateClass)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}
	app.Listen(":" + port)
}

func migrateDatabase(db *gorm.DB) error {
	return db.Exec(schemaSQL).Error
}
