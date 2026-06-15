package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/models"
)

type DepartmentHandler struct {
	OData *client.ODataClient
}

func NewDepartmentHandler(odataClient *client.ODataClient) *DepartmentHandler {
	return &DepartmentHandler{OData: odataClient}
}

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

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	var payload models.DepartmentUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("DEBUG: Ошибка декодирования JSON: %v\n", err)
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
