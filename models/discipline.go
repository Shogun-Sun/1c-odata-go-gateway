package models

// Discipline — структура для фронтенда
type Discipline struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ODataDiscipline — структура для получения из 1С
type ODataDiscipline struct {
	RefKey      string `json:"Ref_Key"`
	Description string `json:"Description"`
}

type ODataDisciplineResponse struct {
	Value []ODataDiscipline `json:"value"`
}

// Payload для создания и изменения
type DisciplinePayload struct {
	Name string `json:"name"`
}

// ODataDisciplineCreate/Update — для отправки в 1С
type ODataDisciplineCreateUpdate struct {
	Description string `json:"Description"`
}
