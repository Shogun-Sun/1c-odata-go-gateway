package handlers

import (
	"academic-booking-api/client"
	"academic-booking-api/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type DisciplineHandler struct {
	OData *client.ODataClient
}

func NewDisciplineHandler(odataClient *client.ODataClient) *DisciplineHandler {
	return &DisciplineHandler{OData: odataClient}
}

// GET /api/v1/disciplines
func (h *DisciplineHandler) GetDisciplines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawData, err := h.OData.Get("Catalog_Дисциплины?$format=json&$select=Ref_Key,Description")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataDisciplineResponse
	json.Unmarshal(rawData, &odataResp)

	disciplines := make([]models.Discipline, 0, len(odataResp.Value))
	for _, d := range odataResp.Value {
		disciplines = append(disciplines, models.Discipline{ID: d.RefKey, Name: d.Description})
	}
	json.NewEncoder(w).Encode(disciplines)
}

// POST /api/v1/disciplines
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

// PATCH /api/v1/disciplines/:id
func (h *DisciplineHandler) UpdateDiscipline(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p models.DisciplinePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	body := models.ODataDisciplineCreateUpdate{Description: p.Name}
	err := h.OData.Patch(fmt.Sprintf("Catalog_Дисциплины(guid'%s')", id), body)
	if err != nil {
		http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DELETE /api/v1/disciplines/:id
func (h *DisciplineHandler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.OData.Delete(fmt.Sprintf("Catalog_Дисциплины(guid'%s')", id))
	if err != nil {
		http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}
