package models

// RoomType определяет тип назначения аудитории.
type RoomType string

// Building определяет учебный корпус площадки.
type Building string

// Константы перечислений, строго соответствующие именам объектов в конфигураторе 1С.
const (
	RoomComputer RoomType = "КомпьютерныйКласс"
	RoomLecture  RoomType = "Лекционный"
	RoomLab      RoomType = "Лаборатория"

	BuildingMain   Building = "ПерваяПлощадка"
	BuildingSecond Building = "ВтораяПлощадка"
	BuildingSormov Building = "Общежитие"
)

// IsValid проверяет, входит ли переданный корпус в список разрешенных.
func (b Building) IsValid() bool {
	switch b {
	case BuildingMain, BuildingSecond, BuildingSormov:
		return true
	}
	return false
}

// IsValid проверяет, валиден ли переданный тип кабинета.
func (rt RoomType) IsValid() bool {
	switch rt {
	case RoomComputer, RoomLecture, RoomLab:
		return true
	}
	return false
}

// Classroom описывает модель аудитории, возвращаемую клиенту (Web API).
type Classroom struct {
	ID       string   `json:"id"`
	Number   string   `json:"number"`
	Capacity int      `json:"capacity"`
	RoomType RoomType `json:"room_type"`
	Building Building `json:"building"`
}

// ODataClassroom описывает структуру аудитории в формате OData 1С.
type ODataClassroom struct {
	RefKey      string   `json:"Ref_Key"`
	Description string   `json:"Description"`
	Capacity    int      `json:"Вместимость,string"` // 1С возвращает числа как строки, тег ",string" чинит это
	RoomType    RoomType `json:"ТипКабинета"`
	Building    Building `json:"Корпус"`
}

// ODataClassroomResponse представляет контейнер верхнего уровня для списка аудиторий из 1С.
type ODataClassroomResponse struct {
	Value []ODataClassroom `json:"value"`
}

// ClassroomCreatePayload содержит данные от фронтенда для создания аудитории.
type ClassroomCreatePayload struct {
	Number   string   `json:"number"`
	Capacity int      `json:"capacity"`
	RoomType RoomType `json:"room_type"`
	Building Building `json:"building"`
}

// ODataClassroomCreate определяет структуру POST-запроса для создания аудитории в 1С.
type ODataClassroomCreate struct {
	Description string   `json:"Description"`
	Capacity    int      `json:"Вместимость"`
	RoomType    RoomType `json:"ТипКабинета"`
	Building    Building `json:"Корпус"`
}

// ClassroomUpdatePayload содержит поля для частичного изменения аудитории (PATCH).
type ClassroomUpdatePayload struct {
	Number   *string   `json:"number,omitempty"`
	Capacity *int      `json:"capacity,omitempty"`
	RoomType *RoomType `json:"room_type,omitempty"`
	Building *Building `json:"building,omitempty"`
}

// ODataClassroomUpdate определяет структуру PATCH-запроса для обновления аудитории в 1С.
type ODataClassroomUpdate struct {
	Description string   `json:"Description,omitempty"`
	Capacity    int      `json:"Вместимость,omitempty"`
	RoomType    RoomType `json:"ТипКабинета,omitempty"`
	Building    Building `json:"Корпус,omitempty"`
}
