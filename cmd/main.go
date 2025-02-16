package main

import (
	"api/internal/repository/memory"
	"api/internal/usecase"
	"log"
	"net/http"

	httpDelivery "api/internal/delivery/http"
)

func main() {
	// User
	repo := memory.NewUserRepositoryMemory()
	userUC := usecase.NewUserUseCase(repo)
	userHandler := httpDelivery.NewUserHandler(userUC)

	// Настраиваем маршруты.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/get", userHandler.GetUser)
	mux.HandleFunc("GET /user/list", userHandler.List)
	mux.HandleFunc("POST /user/create", userHandler.CreateUser)
	mux.HandleFunc("DELETE /user/delete", userHandler.DeleteUser)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
