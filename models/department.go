package models

// Department описывает модель кафедры, возвращаемую клиенту (Web API).
type Department struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	HeadTeacher string `json:"head_teacher"`
}

// ODataDepartment описывает структуру кафедры в формате OData 1С.
type ODataDepartment struct {
	RefKey         string `json:"Ref_Key"`
	Description    string `json:"Description"`
	HeadTeacherKey string `json:"Заведующий_Key"`
}

// ODataDepartmentResponse представляет контейнер верхнего уровня для списка кафедр из 1С.
type ODataDepartmentResponse struct {
	Value []ODataDepartment `json:"value"`
}

// DepartmentCreatePayload содержит данные от фронтенда для создания новой кафедры.
type DepartmentCreatePayload struct {
	Name        string `json:"name"`
	HeadTeacher string `json:"head_teacher"`
}

// ODataDepartmentCreate определяет структуру POST-запроса для создания кафедры в 1С.
type ODataDepartmentCreate struct {
	Description    string `json:"Description"`
	HeadTeacherKey string `json:"Заведующий_Key"`
}

// DepartmentUpdatePayload содержит поля для частичного изменения кафедры (PATCH).
type DepartmentUpdatePayload struct {
	Name        *string `json:"name,omitempty"`
	HeadTeacher *string `json:"head_teacher,omitempty"`
}

// ODataDepartmentUpdate определяет структуру PATCH-запроса для обновления кафедры в 1С.
type ODataDepartmentUpdate struct {
	Description    string `json:"Description,omitempty"`
	HeadTeacherKey string `json:"Заведующий_Key,omitempty"`
}
