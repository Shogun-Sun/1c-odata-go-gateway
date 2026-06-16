package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/models"
)

// TeacherHandler отвечает за обработку HTTP-запросов для сущности «Преподаватели».
type TeacherHandler struct {
	OData *client.ODataClient
}

// NewTeacherHandler возвращает новый экземпляр TeacherHandler.
func NewTeacherHandler(odataClient *client.ODataClient) *TeacherHandler {
	return &TeacherHandler{OData: odataClient}
}

// GetTeachers обрабатывает запрос GET /api/v1/teachers для получения списка всех преподавателей.
func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rawData, err := h.OData.Get("Catalog_Преподаватели?$format=json&$select=Ref_Key,Description,Кафедра_Key,Должность")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataTeacherResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга: "+err.Error(), http.StatusInternalServerError)
		return
	}

	teachers := make([]models.Teacher, 0, len(odataResp.Value))
	for _, ot := range odataResp.Value {
		teachers = append(teachers, models.Teacher{
			ID:         ot.RefKey,
			FullName:   ot.Description,
			Department: ot.DepartmentKey,
			Position:   ot.Position,
		})
	}

	json.NewEncoder(w).Encode(teachers)
}

// CreateTeacher обрабатывает запрос POST /api/v1/teachers для создания нового преподавателя.
func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var payload models.TeacherCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if payload.FullName == "" {
		http.Error(w, "ФИО преподавателя не может быть пустым", http.StatusBadRequest)
		return
	}
	if !payload.Position.IsValid() {
		http.Error(w, fmt.Sprintf("Некорректная должность: %s", payload.Position), http.StatusBadRequest)
		return
	}

	odataBody := models.ODataTeacherCreate{
		Description:   payload.FullName,
		DepartmentKey: payload.Department,
		Position:      payload.Position,
	}

	rawData, err := h.OData.Post("Catalog_Преподаватели?$format=json", odataBody)
	if err != nil {
		http.Error(w, "Ошибка создания в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(rawData)
}

// UpdateTeacher обрабатывает запрос PATCH /api/v1/teachers/{id} для частичного обновления преподавателя.
func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID преподавателя не указан", http.StatusBadRequest)
		return
	}

	var payload models.TeacherUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	odataBody := models.ODataTeacherUpdate{}

	if payload.FullName != nil {
		odataBody.Description = *payload.FullName
	}
	if payload.Department != nil {
		odataBody.DepartmentKey = *payload.Department
	}
	if payload.Position != nil {
		if !(*payload.Position).IsValid() {
			http.Error(w, fmt.Sprintf("Некорректная должность: %s", *payload.Position), http.StatusBadRequest)
			return
		}
		odataBody.Position = *payload.Position
	}

	endpoint := fmt.Sprintf("Catalog_Преподаватели(guid'%s')", id)

	err := h.OData.Patch(endpoint, odataBody)
	if err != nil {
		http.Error(w, "Ошибка обновления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DeleteTeacher обрабатывает запрос DELETE /api/v1/teachers/{id} для удаления преподавателя.
func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID преподавателя не указан", http.StatusBadRequest)
		return
	}

	endpoint := fmt.Sprintf("Catalog_Преподаватели(guid'%s')", id)

	err := h.OData.Delete(endpoint)
	if err != nil {
		http.Error(w, "Ошибка удаления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}
