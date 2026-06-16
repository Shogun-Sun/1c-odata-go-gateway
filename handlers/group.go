package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/models"
)

// GroupHandler отвечает за обработку HTTP-запросов для сущности «Учебные группы».
type GroupHandler struct {
	OData *client.ODataClient
}

// NewGroupHandler возвращает новый экземпляр GroupHandler.
func NewGroupHandler(odataClient *client.ODataClient) *GroupHandler {
	return &GroupHandler{OData: odataClient}
}

// GetGroups обрабатывает запрос GET /api/v1/groups для получения списка всех учебных групп.
func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rawData, err := h.OData.Get("Catalog_УчебныеГруппы?$format=json&$select=Ref_Key,Description,Численность")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataGroupResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга данных 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	frontendGroups := make([]models.Group, 0, len(odataResp.Value))
	for _, og := range odataResp.Value {
		frontendGroups = append(frontendGroups, models.Group{
			ID:       og.RefKey,
			Name:     og.Description,
			Quantity: og.Quantity,
		})
	}

	json.NewEncoder(w).Encode(frontendGroups)
}

// CreateGroup обрабатывает запрос POST /api/v1/groups для создания новой учебной группы.
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var payload models.GroupCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON в запросе", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		http.Error(w, "Имя группы не может быть пустым", http.StatusBadRequest)
		return
	}

	odataBody := models.ODataGroupCreate{
		Description: payload.Name,
		Quantity:    payload.Quantity,
	}

	rawData, err := h.OData.Post("Catalog_УчебныеГруппы?$format=json", odataBody)
	if err != nil {
		http.Error(w, "Ошибка создания группы в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(rawData)
}

// UpdateGroup обрабатывает запрос PATCH /api/v1/groups/{id} для частичного обновления учебной группы.
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID группы не указан", http.StatusBadRequest)
		return
	}

	var payload models.GroupUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	odataBody := models.ODataGroupUpdate{}
	if payload.Name != nil {
		odataBody.Description = *payload.Name
	}
	if payload.Quantity != nil {
		odataBody.Quantity = *payload.Quantity
	}

	endpoint := fmt.Sprintf("Catalog_УчебныеГруппы(guid'%s')", id)

	err := h.OData.Patch(endpoint, odataBody)
	if err != nil {
		http.Error(w, "Ошибка обновления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DeleteGroup обрабатывает запрос DELETE /api/v1/groups/{id} для удаления учебной группы.
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID группы не указан", http.StatusBadRequest)
		return
	}

	endpoint := fmt.Sprintf("Catalog_УчебныеГруппы(guid'%s')", id)

	err := h.OData.Delete(endpoint)
	if err != nil {
		http.Error(w, "Ошибка удаления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}
