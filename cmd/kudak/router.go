package main

import (
	"Kudak/internal/handlers"
	"github.com/gorilla/mux"
)

func InitRouters(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/signup", handler.SignUpHandler).Methods("POST")

	authRote := r.PathPrefix("/signup").Subrouter()

	authRote.HandleFunc("/kindergarten", handler.CreateKindergartenHandler).Methods("POST")
	r.HandleFunc("/kindergarten", handler.GetAllKindergartensHandler).Methods("GET")
	r.HandleFunc("/kindergarten/{id}", handler.GetKindergartenByIDHandler).Methods("GET")

	r.HandleFunc("/education-ministry", handler.CreateEducationMinistryHandler).Methods("POST")
	r.HandleFunc("/education-ministry/{id:[0-9]+}", handler.GetEducationMinistryByIDHandler).Methods("GET")
	r.HandleFunc("/education-ministry", handler.GetAllEducationMinistriesHandler).Methods("GET")
	r.HandleFunc("/education-ministry", handler.UpdateEducationMinistryHandler).Methods("PUT")
	r.HandleFunc("/education-ministry/{id:[0-9]+}", handler.DeleteEducationMinistryHandler).Methods("DELETE")

	// MainDepartment Routes
	r.HandleFunc("/main-department", handler.CreateMainDepartmentHandler).Methods("POST")
	r.HandleFunc("/main-department/{id:[0-9]+}", handler.GetMainDepartmentByIDHandler).Methods("GET")
	r.HandleFunc("/main-department", handler.GetAllMainDepartmentsHandler).Methods("GET")
	r.HandleFunc("/main-department", handler.UpdateMainDepartmentHandler).Methods("PUT")
	r.HandleFunc("/main-department/{id:[0-9]+}", handler.DeleteMainDepartmentHandler).Methods("DELETE")

	return r
}
