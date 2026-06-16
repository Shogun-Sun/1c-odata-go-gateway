package models

// Booking описывает модель бронирования, возвращаемую клиенту (Web API).
type Booking struct {
	ID           string     `json:"id"`
	TeacherID    string     `json:"teacher_id"`
	ClassroomID  string     `json:"classroom_id"`
	DisciplineID string     `json:"discipline_id"`
	GroupID      string     `json:"group_id"`
	StartTime    string     `json:"start_time"`
	EndTime      string     `json:"end_time"`
	Type         LessonType `json:"type"`
	IsPosted     bool       `json:"is_posted"`
}

// BookingCreatePayload содержит данные от фронтенда для создания бронирования.
type BookingCreatePayload struct {
	TeacherID    string     `json:"teacher_id"`
	ClassroomId  string     `json:"classroom_id"`
	DisciplineID string     `json:"discipline_id"`
	GroupID      string     `json:"group_id"`
	StartTime    string     `json:"start_time"` // Формат: YYYY-MM-DDTHH:MM:SS
	EndTime      string     `json:"end_time"`   // Формат: YYYY-MM-DDTHH:MM:SS
	Type         LessonType `json:"type"`
}

// BookingUpdatePayload содержит поля для частичного изменения бронирования (PATCH).
type BookingUpdatePayload struct {
	TeacherID    *string     `json:"teacher_id,omitempty"`
	ClassroomId  *string     `json:"classroom_id,omitempty"`
	DisciplineID *string     `json:"discipline_id,omitempty"`
	GroupID      *string     `json:"group_id,omitempty"`
	StartTime    *string     `json:"start_time,omitempty"`
	EndTime      *string     `json:"end_time,omitempty"`
	Type         *LessonType `json:"type,omitempty"`
}

// ODataBookingRead описывает структуру чтения данных бронирования из OData 1С.
type ODataBookingRead struct {
	RefKey        string     `json:"Ref_Key"`
	Posted        bool       `json:"Posted"`
	TeacherKey    string     `json:"Преподаватель_Key"`
	ClassroomKey  string     `json:"Кабинет_Key"`
	DisciplineKey string     `json:"Дисциплина_Key"`
	GroupKey      string     `json:"Группа_Key"`
	StartTime     string     `json:"ДатаНачала"`
	EndTime       string     `json:"ДатаОкончания"`
	Type          LessonType `json:"ВидЗанятия"`
}

// ODataBookingResponse представляет контейнер верхнего уровня для списка бронирований из 1С.
type ODataBookingResponse struct {
	Value []ODataBookingRead `json:"value"`
}

// ODataBookingCreateUpdate определяет структуру запроса для создания или обновления документа в 1С.
type ODataBookingCreateUpdate struct {
	Date          string     `json:"Date,omitempty"` // Системная дата документа
	Posted        bool       `json:"Posted"`         // Передаем true для вызова ОбработкаПроведения
	TeacherKey    string     `json:"Преподаватель_Key,omitempty"`
	ClassroomKey  string     `json:"Кабинет_Key,omitempty"`
	DisciplineKey string     `json:"Дисциплина_Key,omitempty"`
	GroupKey      string     `json:"Группа_Key,omitempty"`
	StartTime     string     `json:"ДатаНачала,omitempty"`
	EndTime       string     `json:"ДатаОкончания,omitempty"`
	Type          LessonType `json:"ВидЗанятия,omitempty"`
}

// ODataErrorResponse описывает стандартную структуру ошибки, возвращаемую API 1С.
type ODataErrorResponse struct {
	Error struct {
		Message struct {
			Value string `json:"value"`
		} `json:"message"`
	} `json:"error"`
}
