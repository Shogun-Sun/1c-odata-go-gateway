package models

// Department — чистый объект кафедры для фронтенда
type Department struct {
	ID          string `json:"id"`           // UUID кафедры
	Name        string `json:"name"`         // Наименование (Description)
	HeadTeacher string `json:"head_teacher"` // UUID преподавателя-заведующего
}

// ODataDepartment соответствует ответу от 1С
type ODataDepartment struct {
	RefKey         string `json:"Ref_Key"`
	Description    string `json:"Description"`
	HeadTeacherKey string `json:"Заведующий_Key"`
}

type ODataDepartmentResponse struct {
	Value []ODataDepartment `json:"value"`
}

// DepartmentCreatePayload — данные от фронта для создания
type DepartmentCreatePayload struct {
	Name        string `json:"name"`
	HeadTeacher string `json:"head_teacher"`
}

// ODataDepartmentCreate — структура для POST-запроса в 1С
type ODataDepartmentCreate struct {
	Description    string `json:"Description"`
	HeadTeacherKey string `json:"Заведующий_Key"`
}

// DepartmentUpdatePayload — для PATCH
type DepartmentUpdatePayload struct {
	Name        *string `json:"name,omitempty"`
	HeadTeacher *string `json:"head_teacher,omitempty"`
}

// ODataDepartmentUpdate — для PATCH-запроса в 1С
type ODataDepartmentUpdate struct {
	Description    string `json:"Description,omitempty"`
	HeadTeacherKey string `json:"Заведующий_Key,omitempty"`
}
