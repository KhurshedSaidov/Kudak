package handlers

//
//import (
//	"Kudak/internal/service"
//	"Kudak/models"
//	"encoding/base64"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/mux"
//	"io/ioutil"
//	"net/http"
//	"path/filepath"
//	"strconv"
//)
//
//type Handler struct {
//	Service *service.Service
//}
//
//func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
//	var newUser models.User
//	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err := h.Service.CreateUser(newUser.Username, newUser.Password)
//	if err == nil {
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("Пользователь успешно зарегистрирован"))
//	} else if err == service.ErrUsernameTaken {
//		http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
//		return
//	} else {
//		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("Пользователь успешно зарегистрирован"))
//}
//func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
//	var user models.User
//
//	// Декодируем JSON в структуру User
//	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// Проверяем авторизацию
//	authenticated, err := h.Service.Authenticiate(user.Username, user.Password)
//	if err != nil {
//		w.WriteHeader(http.StatusUnauthorized)
//		response := map[string]interface{}{
//			"result": false,
//		}
//		json.NewEncoder(w).Encode(response)
//		return
//	}
//
//	// Отправляем ответ с результатом
//	w.WriteHeader(http.StatusOK)
//	response := map[string]interface{}{
//		"result": authenticated,
//	}
//	json.NewEncoder(w).Encode(response)
//}
//
//// AddKindergartenHandler handles adding a kindergarten
//func (h *Handler) AddKindergartenHandler(w http.ResponseWriter, r *http.Request) {
//	var kindergarten models.Kindergarten
//
//	// Parse multipart form data
//	r.ParseMultipartForm(10 << 20) // Limit the size to 10MB
//
//	// Extract form fields
//	kindergarten.Name = r.FormValue("name")
//	inn, err := strconv.Atoi(r.FormValue("inn"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Print("name or inn error")
//		return
//	}
//	kindergarten.Inn = inn
//	kindergarten.Address = r.FormValue("address")
//	number, err := strconv.Atoi(r.FormValue("number"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Print("Address or number error")
//		return
//	}
//	kindergarten.Number = number
//	kindergarten.Description = r.FormValue("description")
//
//	// Handle the picture file
//	pictureBase64 := r.FormValue("picture")
//	if pictureBase64 != "" {
//		imageFilePath := filepath.Join("uploads", "kindergarten_"+strconv.Itoa(int(kindergarten.ID))+".jpg")
//		if err := decodeBase64Image(pictureBase64, imageFilePath); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			fmt.Print("picture error")
//			return
//		}
//		kindergarten.Picture = imageFilePath
//	}
//
//	// Create kindergarten using service
//	if err := h.Service.CreateKindergarten(&kindergarten); err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	// Response
//	response := map[string]interface{}{
//		"name":        map[string]string{"type": "string", "value": kindergarten.Name},
//		"inn":         map[string]interface{}{"type": "int", "value": kindergarten.Inn},
//		"address":     map[string]string{"type": "string", "value": kindergarten.Address},
//		"number":      map[string]interface{}{"type": "int", "value": kindergarten.Number},
//		"picture":     map[string]string{"type": "string", "value": kindergarten.Picture},
//		"description": map[string]string{"type": "string", "value": kindergarten.Description},
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(response)
//}
//
//// decodeBase64Image decodes a base64 string and saves it as an image file
//func decodeBase64Image(base64Str, filepath string) error {
//	data, err := base64.StdEncoding.DecodeString(base64Str)
//	if err != nil {
//		return err
//	}
//	return ioutil.WriteFile(filepath, data, 0644)
//}
//
//func (h *Handler) GetAllKindergartensHandler(w http.ResponseWriter, r *http.Request) {
//	kindergartens, err := h.Service.GetAllKindergartens()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Ошибка при получении данных"))
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	response := map[string]interface{}{
//		"name":        kindergartens.Name,
//		"inn":         kindergartens.Inn,
//		"address":     kindergartens.Address,
//		"number":      kindergartens.Number,
//		"description": kindergartens.Description,
//	}
//	json.NewEncoder(w).Encode(response)
//}
//
//func (h *Handler) GetKindergartenByIDHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	idStr, ok := vars["id"]
//	if !ok {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Не указан id"))
//		return
//	}
//
//	id, err := strconv.ParseUint(idStr, 10, 32)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Некорректный id"))
//		return
//	}
//
//	kindergarten, err := h.Service.GetKindergartenByID(uint(id))
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Ошибка при получении данных"))
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	response := map[string]interface{}{
//		"name":    kindergarten.Name,
//		"picture": kindergarten.Picture,
//	}
//	json.NewEncoder(w).Encode(response)
//}
