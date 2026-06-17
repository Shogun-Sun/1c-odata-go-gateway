package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/internal/client"
	"academic-booking-api/internal/models"
)

// DepartmentHandler отвечает за обработку HTTP-запросов для сущности «Кафедры».
type DepartmentHandler struct {
	OData *client.ODataClient
}

// NewDepartmentHandler возвращает новый экземпляр DepartmentHandler.
func NewDepartmentHandler(odataClient *client.ODataClient) *DepartmentHandler {
	return &DepartmentHandler{OData: odataClient}
}

// GetDepartments обрабатывает запрос GET /api/v1/departments для получения списка всех кафедр.
func (h *DepartmentHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rawData, err := h.OData.Get("Catalog_Кафедры?$format=json&$select=Ref_Key,Description,Заведующий_Key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataDepartmentResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга: "+err.Error(), http.StatusInternalServerError)
		return
	}

	departments := make([]models.Department, 0, len(odataResp.Value))
	for _, od := range odataResp.Value {
		departments = append(departments, models.Department{
			ID:          od.RefKey,
			Name:        od.Description,
			HeadTeacher: od.HeadTeacherKey,
		})
	}

	json.NewEncoder(w).Encode(departments)
}

// CreateDepartment обрабатывает запрос POST /api/v1/departments для создания новой кафедры.
func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.DepartmentCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	odataBody := models.ODataDepartmentCreate{
		Description:    payload.Name,
		HeadTeacherKey: payload.HeadTeacher,
	}

	rawData, err := h.OData.Post("Catalog_Кафедры?$format=json", odataBody)
	if err != nil {
		http.Error(w, "Ошибка создания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(rawData)
}

// UpdateDepartment обрабатывает запрос PATCH /api/v1/departments/{id} для частичного обновления кафедры.
func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID кафедры не указан", http.StatusBadRequest)
		return
	}

	var payload models.DepartmentUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	odataBody := models.ODataDepartmentUpdate{}
	if payload.Name != nil {
		odataBody.Description = *payload.Name
	}
	if payload.HeadTeacher != nil {
		odataBody.HeadTeacherKey = *payload.HeadTeacher
	}

	endpoint := fmt.Sprintf("Catalog_Кафедры(guid'%s')", id)

	err := h.OData.Patch(endpoint, odataBody)
	if err != nil {
		http.Error(w, "Ошибка обновления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DeleteDepartment обрабатывает запрос DELETE /api/v1/departments/{id} для удаления кафедры.
func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID кафедры не указан", http.StatusBadRequest)
		return
	}

	endpoint := fmt.Sprintf("Catalog_Кафедры(guid'%s')", id)

	err := h.OData.Delete(endpoint)
	if err != nil {
		http.Error(w, "Ошибка удаления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}
