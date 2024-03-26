// main.go

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initDB(host, port, user, password, dbname string) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := pgxpool.New(context.Background(), connectionString)
	return db, err
}

type ServiceI interface {
	myHandler(ctx *fiber.Ctx) error
}
type service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) ServiceI {
	s := service{db: db}
	return &s
}
func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	db, err := initDB(host, port, user, pass, dbname)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	app := fiber.New()
	s := NewService(db)

	app.Get("/user/:id", s.myHandler)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}
	portStr := fmt.Sprintf(":%s", appPort)

	log.Fatal(app.Listen(portStr))
}

func (s *service) myHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var (
		fio     string
		avgBall float32
		cipher  string
	)

	err := s.db.QueryRow(context.Background(), "SELECT fio, avg_ball, cipher FROM users WHERE id=$1", id).Scan(&fio, &avgBall, &cipher)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return ctx.Status(fiber.StatusNotFound).SendString("Пользователь не найден")
		}
		log.Printf("Error querying database: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return ctx.JSON(fiber.Map{
		"fio":      fio,
		"avg_ball": avgBall,
		"cipher":   cipher,
	})
}
