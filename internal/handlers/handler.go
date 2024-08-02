package handlers

import (
	"Kudak/internal/service"
	"Kudak/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
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

func (h *Handler) UploadKindergartenHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form: %v\n", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("picture")
	if err != nil {
		log.Printf("Error getting the file: %v\n", err)
		http.Error(w, "Unable to get the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating the file: %v\n", err)
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error reading the file: %v\n", err)
		http.Error(w, "Unable to read the file", http.StatusInternalServerError)
		return
	}

	_, err = out.Write(fileBytes)
	if err != nil {
		log.Printf("Error writing the file: %v\n", err)
		http.Error(w, "Unable to write the file", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	inn, _ := strconv.Atoi(r.FormValue("inn"))
	phoneNumber := r.FormValue("phoneNumber")
	subtitle := r.FormValue("subtitle")
	latitude, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	longitude, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	address := r.FormValue("address")

	fileRecord := models.Kindergarten{
		Name:        name,
		Path:        filePath,
		Inn:         inn,
		PhoneNumber: phoneNumber,
		Subtitle:    subtitle,
		Latitude:    latitude,
		Longitude:   longitude,
		Address:     address,
	}

	if err := h.Service.CreateKindergarten(&fileRecord); err != nil {
		log.Printf("Error saving file to database: %v\n", err)
		http.Error(w, "Unable to save the file in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fileRecord)
}
func (h *Handler) GetAllKindergartens(w http.ResponseWriter, r *http.Request) {
	kindergartens, err := h.Service.GetAllKindergartens()
	if err != nil {
		log.Printf("Error fetching kindergartens: %v\n", err)
		http.Error(w, "Unable to fetch kindergartens", http.StatusInternalServerError)
		return
	}

	var response []models.KindergartenResponse
	for _, k := range kindergartens {
		response = append(response, models.KindergartenResponse{
			ID:        k.ID,
			Name:      k.Name,
			Inn:       k.Inn,
			Address:   k.Address,
			Number:    k.PhoneNumber,
			CreatedAt: k.CreatedAt,
			Subtitle:  k.Subtitle,
			Latitude:  k.Latitude,
			Longitude: k.Longitude,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetKindergartenBasicInfo(w http.ResponseWriter, r *http.Request) {
	kindergartens, err := h.Service.GetAllKindergartens()
	if err != nil {
		log.Printf("Error fetching kindergartens: %v\n", err)
		http.Error(w, "Unable to fetch kindergartens", http.StatusInternalServerError)
		return
	}

	var response []models.KindergartenBasicInfoResponse
	for _, k := range kindergartens {
		response = append(response, models.KindergartenBasicInfoResponse{
			ID:        k.ID,
			Name:      k.Name,
			Longitude: k.Longitude,
			Latitude:  k.Latitude,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetKindergartenByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid ID: %v\n", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	file, err := h.Service.GetKindergartenByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			log.Printf("Error fetching file: %v\n", err)
			http.Error(w, "Unable to fetch file", http.StatusInternalServerError)
		}
		return
	}

	fileData, err := ioutil.ReadFile(file.Path)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	fileBase64 := base64.StdEncoding.EncodeToString(fileData)
	response := map[string]interface{}{
		"file":    file,
		"picture": fileBase64,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteKindergartenByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid ID: %v\n", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	file, err := h.Service.GetKindergartenByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to fetch data from database", http.StatusInternalServerError)
		return
	}

	if file.Inn == 0 {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	// Удаляем файл из файловой системы
	err = os.Remove(file.Path)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	// Удаляем запись из базы данных
	err = h.Service.DeleteKindergartenByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to delete record from database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "File deleted successfully",
	}

	w.WriteHeader(http.StatusOK)
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
		log.Printf("Ошибка обновления записи: %v\n", err)
		http.Error(w, "Ошибка обновления записи", http.StatusInternalServerError)
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
		log.Printf("Ошибка удаления записи: %v\n", err)
		http.Error(w, "Ошибка удаления записи", http.StatusInternalServerError)
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
		log.Printf("Ошибка обновления записи: %v\n", err)
		http.Error(w, "Ошибка обновления записи", http.StatusInternalServerError)
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
		log.Printf("Ошибка удаления записи: %v\n", err)
		http.Error(w, "Ошибка удаления записи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "MainDepartment deleted successfully"})
}

//
//func (h *Handler) CreateChildHandler(w http.ResponseWriter, r *http.Request) {
//	var child models.Child
//	if err := json.NewDecoder(r.Body).Decode(&child); err != nil {
//		http.Error(w, "Invalid request payload", http.StatusBadRequest)
//		return
//	}
//
//	if err := h.Service.CreateChild(&child); err != nil {
//		http.Error(w, "Failed to create child", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(map[string]interface{}{
//		"message": "Child created successfully",
//		"child":   child,
//	})
//}
//
//func (h *Handler) UpdateAttendancesHandler(w http.ResponseWriter, r *http.Request) {
//	var updateRequest struct {
//		Attendances []models.AttendanceUpdate `json:"attendances"`
//	}
//
//	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
//		http.Error(w, "Invalid request payload", http.StatusBadRequest)
//		return
//	}
//
//	if err := h.Service.UpdateMultipleAttendances(updateRequest.Attendances); err != nil {
//		http.Error(w, "Failed to update attendances", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{"message": "Attendances updated successfully"})
//}
//
//func (h *Handler) GetAllChildrenHandler(w http.ResponseWriter, r *http.Request) {
//	children, err := h.Service.GetAllChildren()
//	if err != nil {
//		http.Error(w, "Failed to get children", http.StatusInternalServerError)
//		return
//	}
//
//	response := make([]map[string]interface{}, len(children))
//	for i, child := range children {
//		response[i] = map[string]interface{}{
//			"id":         child.ID,
//			"full_name":  child.FullName,
//			"birth_date": child.BirthDate,
//			"group":      child.Group,
//		}
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(response)
//}
//
//func (h *Handler) GetLatestAttendanceHandler(w http.ResponseWriter, r *http.Request) {
//	attendances, err := h.Service.GetLatestAttendance()
//	if err != nil {
//		http.Error(w, "Failed to get attendance", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(attendances)
//}
