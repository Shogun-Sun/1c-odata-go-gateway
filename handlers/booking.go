package handlers

import (
	"academic-booking-api/client"
	"academic-booking-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const layout = "2006-01-02T15:04:05"

// BookingHandler отвечает за обработку HTTP-запросов для сущности «Бронирования».
type BookingHandler struct {
	OData *client.ODataClient
}

// NewBookingHandler возвращает новый экземпляр BookingHandler.
func NewBookingHandler(odataClient *client.ODataClient) *BookingHandler {
	return &BookingHandler{OData: odataClient}
}

// GetBookings обрабатывает запрос GET /api/v1/bookings для получения списка всех бронирований.
func (h *BookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rawData, err := h.OData.Get("Document_БронированиеАудитории?$format=json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataBookingResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга данных 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	frontendBookings := make([]models.Booking, 0, len(odataResp.Value))
	for _, ob := range odataResp.Value {
		frontendBookings = append(frontendBookings, models.Booking{
			ID:           ob.RefKey,
			TeacherID:    ob.TeacherKey,
			ClassroomID:  ob.ClassroomKey,
			DisciplineID: ob.DisciplineKey,
			GroupID:      ob.GroupKey,
			StartTime:    ob.StartTime,
			EndTime:      ob.EndTime,
			Type:         ob.Type,
			IsPosted:     ob.Posted,
		})
	}

	json.NewEncoder(w).Encode(frontendBookings)
}

// CreateBooking обрабатывает запрос POST /api/v1/bookings для создания и проведения нового бронирования.
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var payload models.BookingCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON в запросе", http.StatusBadRequest)
		return
	}

	// 1. Парсим строки дат
	startParsed, err := time.ParseInLocation(layout, payload.StartTime, time.Local)
	if err != nil {
		http.Error(w, "Некорректный формат start_time. Используйте YYYY-MM-DDTHH:MM:SS", http.StatusBadRequest)
		return
	}

	endParsed, err := time.ParseInLocation(layout, payload.EndTime, time.Local)
	if err != nil {
		http.Error(w, "Некорректный формат end_time. Используйте YYYY-MM-DDTHH:MM:SS", http.StatusBadRequest)
		return
	}

	// Валидация интервала
	if !endParsed.After(startParsed) {
		http.Error(w, "Дата окончания занятия не может быть меньше или равна дате начала", http.StatusBadRequest)
		return
	}

	if !payload.Type.IsValid() {
		http.Error(w, fmt.Sprintf("Некорректный вид занятия: %s", payload.Type), http.StatusBadRequest)
		return
	}

	formattedStart := startParsed.Format(layout)
	formattedEnd := endParsed.Format(layout)

	odataBody := models.ODataBookingCreateUpdate{
		Date:          formattedStart,
		Posted:        false, // Изначально пишем как черновик
		TeacherKey:    payload.TeacherID,
		ClassroomKey:  payload.ClassroomId,
		DisciplineKey: payload.DisciplineID,
		GroupKey:      payload.GroupID,
		StartTime:     formattedStart,
		EndTime:       formattedEnd,
		Type:          payload.Type,
	}

	// ШАГ 1: Создаем черновик документа в 1С
	rawData, err := h.OData.Post("Document_БронированиеАудитории?$format=json", odataBody)
	if err != nil {
		http.Error(w, "Ошибка создания черновика в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var createdDoc models.ODataBookingRead
	if err := json.Unmarshal(rawData, &createdDoc); err != nil {
		http.Error(w, "Ошибка парсинга ответа создания от 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ШАГ 2: Проводим документ
	postEndpoint := fmt.Sprintf("Document_БронированиеАудитории(guid'%s')/Post", createdDoc.RefKey)
	postData, err := h.OData.Post(postEndpoint, map[string]interface{}{})
	if err != nil {
		var odataErr models.ODataErrorResponse
		if json.Unmarshal([]byte(err.Error()), &odataErr) == nil && odataErr.Error.Message.Value != "" {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf(`{"error": %q}`+"\n", odataErr.Error.Message.Value)))
			return
		}

		http.Error(w, "Ошибка проведения в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(postData)
}

// UpdateBooking обрабатывает запрос PATCH /api/v1/bookings/{id} для изменения существующего бронирования.
func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID бронирования не указан", http.StatusBadRequest)
		return
	}

	var payload models.BookingUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	odataBody := models.ODataBookingCreateUpdate{Posted: true}

	assignStr := func(ptr *string, target *string) {
		if ptr != nil {
			*target = *ptr
		}
	}

	assignStr(payload.TeacherID, &odataBody.TeacherKey)
	assignStr(payload.ClassroomId, &odataBody.ClassroomKey)
	assignStr(payload.DisciplineID, &odataBody.DisciplineKey)
	assignStr(payload.GroupID, &odataBody.GroupKey)

	if payload.StartTime != nil {
		t, err := time.ParseInLocation(layout, *payload.StartTime, time.Local)
		if err != nil {
			http.Error(w, "Некорректный формат start_time", http.StatusBadRequest)
			return
		}
		formatted := t.Format(layout)
		odataBody.StartTime = formatted
		odataBody.Date = formatted
	}
	if payload.EndTime != nil {
		t, err := time.ParseInLocation(layout, *payload.EndTime, time.Local)
		if err != nil {
			http.Error(w, "Некорректный формат end_time", http.StatusBadRequest)
			return
		}
		odataBody.EndTime = t.Format(layout)
	}

	if payload.Type != nil {
		if !(*payload.Type).IsValid() {
			http.Error(w, fmt.Sprintf("Некорректный вид занятия: %s", *payload.Type), http.StatusBadRequest)
			return
		}
		odataBody.Type = *payload.Type
	}

	endpoint := fmt.Sprintf("Document_БронированиеАудитории(guid'%s')", id)
	err := h.OData.Patch(endpoint, odataBody)
	if err != nil {
		var odataErr models.ODataErrorResponse
		if json.Unmarshal([]byte(err.Error()), &odataErr) == nil && odataErr.Error.Message.Value != "" {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf(`{"error": %q}`+"\n", odataErr.Error.Message.Value)))
			return
		}

		http.Error(w, "Ошибка обновления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}` + "\n"))
}

// DeleteBooking обрабатывает запрос DELETE /api/v1/bookings/{id} для удаления бронирования.
func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID бронирования не указан", http.StatusBadRequest)
		return
	}

	endpoint := fmt.Sprintf("Document_БронированиеАудитории(guid'%s')", id)
	if err := h.OData.Delete(endpoint); err != nil {
		http.Error(w, "Ошибка удаления в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}` + "\n"))
}

// DeleteAllBookings обрабатывает запрос DELETE /api/v1/bookings для массового удаления всех бронирований.
func (h *BookingHandler) DeleteAllBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rawData, err := h.OData.Get("Document_БронированиеАудитории?$format=json&$select=Ref_Key")
	if err != nil {
		http.Error(w, "Ошибка получения списка для удаления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var odataResp models.ODataBookingResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга списка 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(odataResp.Value) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","deleted_count":0,"message":"Нет документов для удаления"}` + "\n"))
		return
	}

	deletedCount := 0
	for _, ob := range odataResp.Value {
		endpoint := fmt.Sprintf("Document_БронированиеАудитории(guid'%s')", ob.RefKey)
		if err := h.OData.Delete(endpoint); err != nil {
			continue
		}
		deletedCount++
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"success","deleted_count":%d,"message":"Все документы успешно удалены"}`+"\n", deletedCount)))
}
