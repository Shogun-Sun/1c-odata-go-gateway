package handlers

import (
	"encoding/json"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/models"
)

// GroupHandler группирует методы для работы с группами
type GroupHandler struct {
	OData *client.ODataClient
}

// NewGroupHandler — конструктор хендлера
func NewGroupHandler(odataClient *client.ODataClient) *GroupHandler {
	return &GroupHandler{OData: odataClient}
}

// GetGroups обрабатывает запрос GET /api/v1/groups
func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // CORS на перспективу для фронта

	// 1. Запрашиваем данные у 1С
	// Запрашиваем сразу с параметром $select, чтобы не тянуть лишние системные поля 1С
	rawData, err := h.OData.Get("Catalog_УчебныеГруппы?$format=json&$select=Ref_Key,Description,Численность")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Десериализуем JSON из 1С во внутреннюю структуру
	var odataResp models.ODataGroupResponse
	if err := json.Unmarshal(rawData, &odataResp); err != nil {
		http.Error(w, "Ошибка парсинга данных 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Мапим (пересобираем) «грязный» JSON 1С в чистый формат нашего API
	frontendGroups := make([]models.Group, 0, len(odataResp.Value))
	for _, og := range odataResp.Value {
		frontendGroups = append(frontendGroups, models.Group{
			ID:       og.RefKey,
			Name:     og.Description,
			Quantity: og.Quantity, // Наше числовое поле улетело на фронт
		})
	}

	// 4. Отдаем результат в HTTP-ответ
	json.NewEncoder(w).Encode(frontendGroups)
}

// POST /api/v1/groups
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Декодируем то, что прислал пользователь (фронтенд)
	var payload models.GroupCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный JSON в запросе", http.StatusBadRequest)
		return
	}

	// Валидация на скорую руку
	if payload.Name == "" {
		http.Error(w, "Имя группы не может быть пустым", http.StatusBadRequest)
		return
	}

	// 2. Формируем структуру для 1С
	odataBody := models.ODataGroupCreate{
		Description: payload.Name,
		Quantity:    payload.Quantity,
	}

	// 3. Отправляем POST-запрос в 1С OData
	rawData, err := h.OData.Post("Catalog_УчебныеГруппы?$format=json", odataBody)
	if err != nil {
		http.Error(w, "Ошибка создания группы в 1С: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Отдаем клиенту то, что вернула 1С (со всеми системными полями и сгенерированным Ref_Key)
	w.WriteHeader(http.StatusCreated)
	w.Write(rawData)
}
