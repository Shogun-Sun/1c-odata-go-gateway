package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/models"
)

// DisciplineHandler отвечает за обработку HTTP-запросов для сущности «Дисциплины».
type DisciplineHandler struct {
	OData *client.ODataClient
}

// NewDisciplineHandler возвращает новый экземпляр DisciplineHandler.
func NewDisciplineHandler(odataClient *client.ODataClient) *DisciplineHandler {
	return &DisciplineHandler{OData: odataClient}
}

// GetDisciplines обрабатывает запрос GET /api/v1/disciplines для получения списка всех дисциплин.
func (h *DisciplineHandler) GetDisciplines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawData, err := h.OData.Get("Catalog_Дисциплины?$format=json&$select=Ref_Key,Description")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataDisciplineResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга: "+err.Error(), http.StatusInternalServerError)
		return
	}

	disciplines := make([]models.Discipline, 0, len(odataResp.Value))
	for _, d := range odataResp.Value {
		disciplines = append(disciplines, models.Discipline{
			ID:   d.RefKey,
			Name: d.Description,
		})
	}

	json.NewEncoder(w).Encode(disciplines)
}

// CreateDiscipline обрабатывает запрос POST /api/v1/disciplines для создания новой дисциплины.
func (h *DisciplineHandler) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p models.DisciplinePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	body := models.ODataDisciplineCreateUpdate{Description: p.Name}
	rawData, err := h.OData.Post("Catalog_Дисциплины?$format=json", body)
	if err != nil {
		http.Error(w, "Ошибка создания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(rawData)
}

// UpdateDiscipline обрабатывает запрос PATCH /api/v1/disciplines/{id} для обновления дисциплины.
func (h *DisciplineHandler) UpdateDiscipline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID дисциплины не указан", http.StatusBadRequest)
		return
	}

	var p models.DisciplinePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	body := models.ODataDisciplineCreateUpdate{Description: p.Name}
	endpoint := fmt.Sprintf("Catalog_Дисциплины(guid'%s')", id)

	err := h.OData.Patch(endpoint, body)
	if err != nil {
		http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DeleteDiscipline обрабатывает запрос DELETE /api/v1/disciplines/{id} для удаления дисциплины.
func (h *DisciplineHandler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID дисциплины не указан", http.StatusBadRequest)
		return
	}

	endpoint := fmt.Sprintf("Catalog_Дисциплины(guid'%s')", id)

	err := h.OData.Delete(endpoint)
	if err != nil {
		http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}
