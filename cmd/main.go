package main

import (
	httpHandlers "api/internal/delivery/http/handlers"
	"api/internal/delivery/http/middlewares"
	"api/internal/logger"
	"api/internal/repository/memory"
	"api/internal/usecase"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	logConfig := logger.Config{
		Level:      slog.LevelInfo,
		OutputPath: logger.OutputPathStdOut,
		Format:     logger.FormatText,
	}

	lgr, err := logger.New(logConfig)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: Вынести в конфиг
	connStr := "postgres://user:password@localhost:5432/dbname?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	// User

	//userRepo := postgres.NewPostgresUserRepository(pool)
	userRepo := memory.NewUserRepositoryMemory()
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := httpHandlers.NewUserHandler(userUC)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/get", userHandler.GetUser)
	mux.HandleFunc("GET /user/list", userHandler.List)
	mux.HandleFunc("POST /user/create", userHandler.CreateUser)
	mux.HandleFunc("DELETE /user/delete", userHandler.DeleteUser)

	loggedMux := middlewares.Logging(mux, lgr)

	lgr.Info("Server started on :8080")
	if err = http.ListenAndServe(":8080", loggedMux); err != nil {
		lgr.Error("Fail to run server: ", err.Error())
		panic(err)
	}
}
