package handlers

import (
	"Kudak/internal/service"
	"Kudak/models"
	"encoding/json"
	"net/http"
)

type Handler struct {
	Service *service.Service
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.CreateUser(newUser.Username, newUser.Password)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователь успешно зарегистрирован"))
	} else if err == service.ErrUsernameTaken {
		http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
		return
	} else {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно зарегистрирован"))
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.Service.Authenticiate(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	user2, err := h.Service.Repository.GetUserByUsername(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Успешный вход",
		"user":    user2,
	}
	json.NewEncoder(w).Encode(response)
}
