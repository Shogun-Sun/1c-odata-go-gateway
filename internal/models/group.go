package models

// Group описывает модель учебной группы, возвращаемую клиенту (Web API).
type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// ODataGroup описывает структуру учебной группы в формате OData 1С.
type ODataGroup struct {
	RefKey      string `json:"Ref_Key"`
	Description string `json:"Description"`
	Quantity    int    `json:"Численность,string"` // 1С возвращает числа как строки, тег ",string" чинит это
}

// ODataGroupResponse представляет контейнер верхнего уровня для списка групп из 1С.
type ODataGroupResponse struct {
	Value []ODataGroup `json:"value"`
}

// GroupCreatePayload содержит данные от фронтенда для создания новой учебной группы.
type GroupCreatePayload struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// ODataGroupCreate определяет структуру POST-запроса для создания группы в 1С.
type ODataGroupCreate struct {
	Description string `json:"Description"`
	Quantity    int    `json:"Численность"`
}

// GroupUpdatePayload содержит поля для частичного изменения учебной группы (PATCH).
type GroupUpdatePayload struct {
	Name     *string `json:"name,omitempty"`
	Quantity *int    `json:"quantity,omitempty"`
}

// ODataGroupUpdate определяет структуру PATCH-запроса для обновления группы в 1С.
type ODataGroupUpdate struct {
	Description string `json:"Description,omitempty"`
	Quantity    int    `json:"Численность,omitempty"`
}
