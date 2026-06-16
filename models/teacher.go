package models

// Position определяет должность преподавателя в виде строго типизированной строки.
type Position string

// Константы должностей, соответствующие значениям перечисления в конфигураторе 1С.
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

// Teacher описывает модель преподавателя, возвращаемую клиенту (Web API).
type Teacher struct {
	ID         string   `json:"id"`
	FullName   string   `json:"full_name"`
	Department string   `json:"department"`
	Position   Position `json:"position"`
}

// ODataTeacher описывает структуру преподавателя в формате OData 1С.
type ODataTeacher struct {
	RefKey        string   `json:"Ref_Key"`
	Description   string   `json:"Description"`
	DepartmentKey string   `json:"Кафедра_Key"`
	Position      Position `json:"Должность"`
}

// ODataTeacherResponse представляет контейнер верхнего уровня для списка преподавателей из 1С.
type ODataTeacherResponse struct {
	Value []ODataTeacher `json:"value"`
}

// TeacherCreatePayload содержит данные от фронтенда для создания нового преподавателя.
type TeacherCreatePayload struct {
	FullName   string   `json:"full_name"`
	Department string   `json:"department"`
	Position   Position `json:"position"`
}

// ODataTeacherCreate определяет структуру POST-запроса для создания преподавателя в 1С.
type ODataTeacherCreate struct {
	Description   string   `json:"Description"`
	DepartmentKey string   `json:"Кафедра_Key,omitempty"`
	Position      Position `json:"Должность"`
}

// TeacherUpdatePayload содержит поля для частичного изменения данных преподавателя (PATCH).
type TeacherUpdatePayload struct {
	FullName   *string   `json:"full_name,omitempty"`
	Department *string   `json:"department,omitempty"`
	Position   *Position `json:"position,omitempty"`
}

// ODataTeacherUpdate определяет структуру PATCH-запроса для обновления преподавателя в 1С.
type ODataTeacherUpdate struct {
	Description   string   `json:"Description,omitempty"`
	DepartmentKey string   `json:"Кафедра_Key,omitempty"`
	Position      Position `json:"Должность,omitempty"`
}
