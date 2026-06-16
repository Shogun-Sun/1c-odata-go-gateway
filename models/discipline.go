package models

// Discipline описывает модель учебной дисциплины, возвращаемую клиенту (Web API).
type Discipline struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ODataDiscipline описывает структуру дисциплины в формате OData 1С.
type ODataDiscipline struct {
	RefKey      string `json:"Ref_Key"`
	Description string `json:"Description"`
}

// ODataDisciplineResponse представляет контейнер верхнего уровня для списка дисциплин из 1С.
type ODataDisciplineResponse struct {
	Value []ODataDiscipline `json:"value"`
}

// DisciplinePayload содержит данные от фронтенда для создания или полного изменения дисциплины.
type DisciplinePayload struct {
	Name string `json:"name"`
}

// ODataDisciplineCreateUpdate определяет структуру запроса для создания или обновления дисциплины в 1С.
type ODataDisciplineCreateUpdate struct {
	Description string `json:"Description"`
}
