package main

import (
	"database/sql"
	"fmt"
	"github.com/Fyefhqdishka/web-project/pkg/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	logger := initLogging()

	if err := loadEnv(logger); err != nil {
		logger.Error("Ошибка загрузки окружения", "error", err)
		os.Exit(1)
	}

	db, err := connectToDB()
	if err != nil {
		logger.Error("Ошибка подключения к базе данных", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	r := mux.NewRouter()

	routes.RegisterRoutes(r, db, logger)

	port := ":8000"
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal()
	}
}

func initLogging() *slog.Logger {
	logFileName := "logs/app-" + time.Now().Format("2006-01-02") + ".log"
	logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Не удалось открыть файл для логов", "error", err)
		os.Exit(1)
	}

	handler := slog.NewTextHandler(logfile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}

func connectToDB() (*sql.DB, error) {
	dbPass := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=deadlock sslmode=disable", dbPass)
	return sql.Open("postgres", connStr)
}

func loadEnv(logger *slog.Logger) error {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Ошибка загрузки окружения", "error", err)
		return err
	}
	return nil
}
