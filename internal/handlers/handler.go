package handlers

import (
	"Kudak/internal/service"
	"Kudak/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func generateRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type Handler struct {
	Service *service.Service
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Не удалось декодировать JSON. Проверьте правильность формата.", http.StatusBadRequest)
		log.Printf("Ошибка декодирования JSON: %v", err)
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

	// Декодируем JSON в структуру User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Не удалось декодировать JSON. Проверьте правильность формата.", http.StatusBadRequest)
		log.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	// Проверяем авторизацию
	authenticated, err := h.Service.Authenticiate(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{
			"result": false,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Отправляем ответ с результатом
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"result": authenticated,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateKindergartenHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Name        string   `json:"name"`
		Inn         int      `json:"inn"`
		Address     string   `json:"address"`
		Number      int      `json:"number"`
		Pictures    []string `json:"pictures"`
		Description string   `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Не удалось декодировать JSON. Проверьте правильность формата.", http.StatusBadRequest)
		log.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	var pictureModels []models.KindergartenPicture
	for _, base64Image := range requestData.Pictures {
		imageData, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			http.Error(w, "Не удалось декодировать изображение. Проверьте правильность формата base64.", http.StatusBadRequest)
			log.Printf("Ошибка декодирования base64: %v", err)
			return
		}

		filename := filepath.Join("uploads", "kindergarten_"+generateRandomString(10)+".jpg")
		if err := ioutil.WriteFile(filename, imageData, 0644); err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		pictureModels = append(pictureModels, models.KindergartenPicture{PicturePath: filename})
	}

	kindergarten := &models.Kindergarten{
		Name:        requestData.Name,
		Inn:         requestData.Inn,
		Address:     requestData.Address,
		Number:      requestData.Number,
		Picture:     pictureModels,
		Description: requestData.Description,
	}

	if err := h.Service.CreateKindergarten(kindergarten); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(kindergarten)
}

func decodeBase64Image(base64Str, filePath string) error {
	// Удаляем префикс data URL, если он есть
	if strings.HasPrefix(base64Str, "data:image/") {
		parts := strings.SplitN(base64Str, ",", 2)
		if len(parts) != 2 {
			return errors.New("некорректный формат base64 строки")
		}
		base64Str = parts[1]
	}

	// Декодируем строку base64
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("ошибка декодирования base64: %v", err)
	}

	// Сохраняем изображение в файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}

func (h *Handler) GetAllKindergartensHandler(w http.ResponseWriter, r *http.Request) {
	kindergartens, err := h.Service.GetAllKindergartens()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при получении данных"))
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"name":        kindergartens.Name,
		"inn":         kindergartens.Inn,
		"address":     kindergartens.Address,
		"number":      kindergartens.Number,
		"description": kindergartens.Description,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetKindergartenByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Не указан id в URL параметрах", http.StatusBadRequest)
		log.Println("Не указан id в URL параметрах")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		log.Printf("Некорректный id: %v\n", err)
		return
	}

	kindergarten, err := h.Service.GetKindergartenByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при получении данных"))
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"name":    kindergarten.Name,
		"picture": kindergarten.Picture,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateEducationMinistryHandler(w http.ResponseWriter, r *http.Request) {
	var em models.EducationMinistry
	if err := json.NewDecoder(r.Body).Decode(&em); err != nil {
		log.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Невозможно декодировать JSON", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateEducationMinistry(&em); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "EducationMinistry created successfully"})
}

func (h *Handler) GetAllEducationMinistriesHandler(w http.ResponseWriter, r *http.Request) {
	ministries, err := h.Service.GetAllEducationMinistries()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ministries)
}

func (h *Handler) GetEducationMinistryByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования id: %v\n", err)
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	em, err := h.Service.GetEducationMinistryByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(em)
}

func (h *Handler) UpdateEducationMinistryHandler(w http.ResponseWriter, r *http.Request) {
	var em models.EducationMinistry
	if err := json.NewDecoder(r.Body).Decode(&em); err != nil {
		log.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Невозможно декодировать JSON", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateEducationMinistry(&em); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "EducationMinistry updated successfully"})
}

func (h *Handler) DeleteEducationMinistryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования id: %v\n", err)
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteEducationMinistry(uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "EducationMinistry deleted successfully"})
}

func (h *Handler) CreateMainDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var md models.MainDepartment
	if err := json.NewDecoder(r.Body).Decode(&md); err != nil {
		log.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Невозможно декодировать JSON", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateMainDepartment(&md); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "MainDepartment created successfully"})
}

func (h *Handler) GetAllMainDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	departments, err := h.Service.GetAllMainDepartments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(departments)
}

func (h *Handler) GetMainDepartmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования id: %v\n", err)
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	md, err := h.Service.GetMainDepartmentByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(md)
}

func (h *Handler) UpdateMainDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var md models.MainDepartment
	if err := json.NewDecoder(r.Body).Decode(&md); err != nil {
		log.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Невозможно декодировать JSON", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateMainDepartment(&md); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "MainDepartment updated successfully"})
}

func (h *Handler) DeleteMainDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования id: %v\n", err)
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteMainDepartment(uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "MainDepartment deleted successfully"})
}
