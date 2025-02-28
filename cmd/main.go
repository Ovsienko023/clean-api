package main

import (
	httpHandlers "api/internal/delivery/http/handlers"
	"api/internal/delivery/http/middlewares"
	"api/internal/logger"
	"api/internal/repository/memory"
	"api/internal/usecase"
	"log/slog"
	"net/http"
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

	// User
	repo := memory.NewUserRepositoryMemory()
	userUC := usecase.NewUserUseCase(repo)
	userHandler := httpHandlers.NewUserHandler(userUC)

	// Настраиваем маршруты.
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
