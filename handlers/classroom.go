package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/models"
)

// ClassroomHandler отвечает за обработку HTTP-запросов для кабинетов
type ClassroomHandler struct {
	OData *client.ODataClient // Клиент для связи с 1С
}

// NewClassroomHandler — конструктор для инициализации хендлера
func NewClassroomHandler(odataClient *client.ODataClient) *ClassroomHandler {
	return &ClassroomHandler{OData: odataClient}
}

// GET /api/v1/classrooms — Получение списка всех кабинетов
func (h *ClassroomHandler) GetClassrooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Запрашиваем данные у 1С с нужными полями
	rawData, err := h.OData.Get("Catalog_Кабинеты?$format=json&$select=Ref_Key,Description,Вместимость,ТипКабинета,Корпус")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataClassroomResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Переводим данные из формата 1С в чистый JSON для фронтенда
	classrooms := make([]models.Classroom, 0, len(odataResp.Value))
	for _, oc := range odataResp.Value {
		classrooms = append(classrooms, models.Classroom{
			ID:       oc.RefKey,
			Number:   oc.Description,
			Capacity: oc.Capacity,
			RoomType: oc.RoomType,
			Building: oc.Building,
		})
	}

	json.NewEncoder(w).Encode(classrooms)
}

// POST /api/v1/classrooms — Создание нового кабинета
func (h *ClassroomHandler) CreateClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var payload models.ClassroomCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	// Проверка входящих значений на соответствие перечислениям 1С
	if !payload.Building.IsValid() {
		http.Error(w, fmt.Sprintf("Некорректный корпус: %s", payload.Building), http.StatusBadRequest)
		return
	}
	if !payload.RoomType.IsValid() {
		http.Error(w, fmt.Sprintf("Некорректный тип кабинета: %s", payload.RoomType), http.StatusBadRequest)
		return
	}

	odataBody := models.ODataClassroomCreate{
		Description: payload.Number,
		Capacity:    payload.Capacity,
		RoomType:    payload.RoomType,
		Building:    payload.Building,
	}

	// Отправляем пакет на создание в 1С
	rawData, err := h.OData.Post("Catalog_Кабинеты?$format=json", odataBody)
	if err != nil {
		http.Error(w, "Ошибка создания в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(rawData)
}

// PATCH /api/v1/classrooms/{id} — Частичное обновление данных кабинета
func (h *ClassroomHandler) UpdateClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID кабинета не указан", http.StatusBadRequest)
		return
	}

	var payload models.ClassroomUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	// Заполняем структуру обновления только теми полями, которые пришли в запросе
	odataBody := models.ODataClassroomUpdate{}

	if payload.Number != nil {
		odataBody.Description = *payload.Number
	}
	if payload.Capacity != nil {
		odataBody.Capacity = *payload.Capacity
	}
	if payload.RoomType != nil {
		if !(*payload.RoomType).IsValid() {
			http.Error(w, fmt.Sprintf("Некорректный тип кабинета: %s", *payload.RoomType), http.StatusBadRequest)
			return
		}
		odataBody.RoomType = *payload.RoomType
	}
	if payload.Building != nil {
		if !(*payload.Building).IsValid() {
			http.Error(w, fmt.Sprintf("Некорректный корпус: %s", *payload.Building), http.StatusBadRequest)
			return
		}
		odataBody.Building = *payload.Building
	}

	endpoint := fmt.Sprintf("Catalog_Кабинеты(guid'%s')", id)

	err := h.OData.Patch(endpoint, odataBody)
	if err != nil {
		http.Error(w, "Ошибка обновления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DELETE /api/v1/classrooms/{id} — Удаление кабинета
func (h *ClassroomHandler) DeleteClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID кабинета не указан", http.StatusBadRequest)
		return
	}

	endpoint := fmt.Sprintf("Catalog_Кабинеты(guid'%s')", id)

	err := h.OData.Delete(endpoint)
	if err != nil {
		http.Error(w, "Ошибка удаления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}
