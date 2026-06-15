package models

// Position определяет должность преподавателя (строго типизированная строка).
type Position string

// Перечень разрешенных должностей из перечисления в 1С.
const (
	PositionTeacher Position = "Преподаватель"
)

// IsValid проверяет, входит ли переданная должность в список разрешенных.
func (p Position) IsValid() bool {
	switch p {
	case PositionTeacher:
		return true
	}
	return false
}

// Teacher описывает чистый объект преподавателя для фронтенда (GET)
type Teacher struct {
	ID         string   `json:"id"`         // UUID преподавателя
	FullName   string   `json:"full_name"`  // ФИО (Description в 1С)
	Department string   `json:"department"` // UUID кафедры (Кафедра_Key)
	Position   Position `json:"position"`   // Должность
}

// ODataTeacher соответствует JSON-объекту, который возвращает 1С при GET-запросе
type ODataTeacher struct {
	RefKey        string   `json:"Ref_Key"`     // UUID элемента
	Description   string   `json:"Description"` // ФИО преподавателя
	DepartmentKey string   `json:"Кафедра_Key"` // Ссылка на GUID кафедры
	Position      Position `json:"Должность"`   // Значение перечисления Должности
}

// ODataTeacherResponse описывает обертку верхнего уровня для GET-ответов от OData 1С
type ODataTeacherResponse struct {
	Value []ODataTeacher `json:"value"`
}

// TeacherCreatePayload прилетает от фронтенда при создании (POST)
type TeacherCreatePayload struct {
	FullName   string   `json:"full_name"`
	Department string   `json:"department"` // UUID кафедры
	Position   Position `json:"position"`
}

// ODataTeacherCreate определяет структуру для отправки POST-запроса в 1С
type ODataTeacherCreate struct {
	Description   string   `json:"Description"`
	DepartmentKey string   `json:"Кафедра_Key"` // 1С OData ждет суффикс _Key для ссылочных реквизитов
	Position      Position `json:"Должность"`
}

// TeacherUpdatePayload прилетает от фронтенда для изменения (PATCH)
type TeacherUpdatePayload struct {
	FullName   *string   `json:"full_name,omitempty"`
	Department *string   `json:"department,omitempty"`
	Position   *Position `json:"position,omitempty"`
}

// ODataTeacherUpdate передает измененные реквизиты в 1С
type ODataTeacherUpdate struct {
	Description   string   `json:"Description,omitempty"`
	DepartmentKey string   `json:"Кафедра_Key,omitempty"`
	Position      Position `json:"Должность,omitempty"`
}
