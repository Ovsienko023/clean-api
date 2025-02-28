package handlers

import (
	http2 "api/internal/delivery/http"
	"api/internal/domain"
	"api/internal/usecase"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// UserHandler отвечает за обработку HTTP-запросов для пользователей.
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	validate    *validator.Validate
}

// NewUserHandler создаёт новый экземпляр UserHandler.
func NewUserHandler(u *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: u,
		validate:    validator.New(),
	}
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// GetUser обрабатывает получение пользователя по ID.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	errorContainer := http2.ErrorResponse{}

	request := &GetUserRequest{
		ID: r.URL.Query().Get("id"),
	}

	// Валидация структуры запроса
	if err := h.validate.Struct(request); err != nil {
		errorContainer.Done(w, err)
		return
	}

	user, err := h.userUseCase.GetUser(r.Context(), request.ID)
	if err != nil {
		errorContainer.Done(w, err)
		return
	}

	http2.JsonResponse(w, http.StatusOK, user)
}

type ListUserRequest struct {
	Limit  *uint32 `json:"limit"`
	Offset *uint32 `json:"offset"`
}

// List возвращает список всех пользователей.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	errorContainer := http2.ErrorResponse{}

	request := domain.QueryOptions{}

	users, err := h.userUseCase.ListUser(r.Context(), request)
	if err != nil {
		errorContainer.Done(w, err)
		return
	}

	http2.JsonResponse(w, http.StatusOK, users)
}

// CreateUserRequest представляет входящие данные для создания пользователя.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

// CreateUser обрабатывает создание нового пользователя.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	errorContainer := http2.ErrorResponse{}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация структуры запроса
	if err := h.validate.Struct(req); err != nil {
		errorContainer.Done(w, err)
		return
	}

	// Преобразуем CreateUserRequest в доменную модель User.
	user := domain.User{
		Name: req.Name,
	}

	id, err := h.userUseCase.CreateUser(r.Context(), user)
	if err != nil {
		errorContainer.Done(w, err)
		return
	}

	response := struct {
		ID *string `json:"id"`
	}{
		ID: id,
	}

	http2.JsonResponse(w, http.StatusCreated, response)
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteUser обрабатывает удаление пользователя по ID.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	errorContainer := http2.ErrorResponse{}

	request := DeleteUserRequest{
		ID: r.URL.Query().Get("id"),
	}

	// Валидация структуры запроса
	if err := h.validate.Struct(request); err != nil {
		errorContainer.Done(w, err)
		return
	}

	if err := h.userUseCase.DeleteUser(r.Context(), request.ID); err != nil {
		errorContainer.Done(w, err)
		return
	}

	http2.JsonResponse(w, http.StatusNoContent, nil)
}
