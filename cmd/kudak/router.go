package main

import (
	"Kudak/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRouters(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/signup", handler.SignUpHandler).Methods("POST")

	//authRote := r.PathPrefix("/signup").Subrouter()

	r.HandleFunc("/upload", handler.UploadKindergartenHandler).Methods("POST")
	r.HandleFunc("/files", handler.GetAllKindergartens).Methods("GET")
	r.HandleFunc("/files/{id}", handler.GetKindergartenByIDHandler).Methods("GET")
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

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

	//api := r.PathPrefix("/api").Subrouter()
	//
	//api.HandleFunc("/children", handler.CreateChildHandler).Methods("POST")
	//api.HandleFunc("/children", handler.GetAllChildrenHandler).Methods("GET")
	//api.HandleFunc("/attendance", handler.UpdateAttendancesHandler).Methods("POST")
	//api.HandleFunc("/attendance/latest", handler.GetLatestAttendanceHandler).Methods("GET")

	return r
}
