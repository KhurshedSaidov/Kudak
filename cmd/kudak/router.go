package main

import (
	"Kudak/internal/handlers"
	"github.com/gorilla/mux"
)

func InitRouters(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/signup", handler.SignUpHandler).Methods("POST")

	return r
}
