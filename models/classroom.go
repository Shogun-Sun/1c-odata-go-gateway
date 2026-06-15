package models

// RoomType определяет тип кабинета (строго типизированная строка).
// Нужен для исключения опечаток при передаче типов комнат в коде.
type RoomType string

// Building определяет учебный корпус или площадку.
// Предотвращает сквозной проброс некорректных текстовых названий площадок.
type Building string

// Фиксированный перечень разрешенных значений (аналог Enum в 1С).
// Значения констант должны строго совпадать с ИМЕНАМИ объектов в конфигураторе 1С.
const (
	RoomComputer RoomType = "КомпьютерныйКласс" // Кабинет с ПК
	RoomLecture  RoomType = "Лекционный"        // Амфитеатр или поточная аудитория
	RoomLab      RoomType = "Лаборатория"       // Специализированная лаборатория

	BuildingMain   Building = "ПерваяПлощадка" // Главный корпус
	BuildingSecond Building = "ВтораяПлощадка" // Дополнительный учебный корпус
	BuildingSormov Building = "Общежитие"      // Помещения в здании общежития
)

// IsValid проверяет, входит ли переданный корпус в список разрешенных в 1С.
// Используется в хендлерах как «щит» перед отправкой запроса к базе данных.
func (b Building) IsValid() bool {
	switch b {
	case BuildingMain, BuildingSecond, BuildingSormov:
		return true
	}
	return false
}

// IsValid проверяет, валиден ли переданный тип кабинета.
// Гарантирует, что фронтенд прислал известное системе назначение комнаты.
func (rt RoomType) IsValid() bool {
	switch rt {
	case RoomComputer, RoomLecture, RoomLab:
		return true
	}
	return false
}

// Classroom описывает чистый объект кабинета, который отдается на фронтенд.
// Поля приведены к стандартному веб-виду (id вместо Ref_Key, room_type вместо ТипКабинета).
type Classroom struct {
	ID       string   `json:"id"`        // Уникальный UUID объекта
	Number   string   `json:"number"`    // Номер кабинета (Description в 1С)
	Capacity int      `json:"capacity"`  // Вместимость мест
	RoomType RoomType `json:"room_type"` // Тип (Компьютерный, Лекционный...)
	Building Building `json:"building"`  // Учебный корпус
}

// ODataClassroom соответствует сырому JSON-объекту, который возвращает 1С при GET-запросе.
// Необходим для корректного демаршалинга (парсинга) кириллических полей 1С.
type ODataClassroom struct {
	RefKey      string   `json:"Ref_Key"`            // Внутренний UUID элемента в 1С
	Description string   `json:"Description"`        // Наименование (номер кабинета)
	Capacity    int      `json:"Вместимость,string"` // 1С часто возвращает числа как строки, тег ",string" автоматически конвертирует их в int
	RoomType    RoomType `json:"ТипКабинета"`        // Реквизит перечисления из 1С
	Building    Building `json:"Корпус"`             // Реквизит перечисления из 1С
}

// ODataClassroomResponse описывает обертку верхнего уровня для GET-ответов от OData 1С.
// В OData все массивы объектов всегда прилетают внутри JSON-поля "value".
type ODataClassroomResponse struct {
	Value []ODataClassroom `json:"value"` // Массив сырых объектов из 1С
}

// ClassroomCreatePayload описывает структуру JSON, которую присылает фронтенд при POST-запросе.
// Здесь нет поля ID, так как идентификатор генерируется на стороне 1С при создании.
type ClassroomCreatePayload struct {
	Number   string   `json:"number"`    // Номер создаваемого кабинета
	Capacity int      `json:"capacity"`  // Вместимость мест
	RoomType RoomType `json:"room_type"` // Назначение комнаты
	Building Building `json:"building"`  // В каком корпусе находится
}

// ODataClassroomCreate определяет структуру данных для отправки POST-запроса непосредственно в OData 1С.
// Переводит названия полей на кириллицу, понятную конфигурации 1С.
type ODataClassroomCreate struct {
	Description string   `json:"Description"` // Передаем номер в стандартное Наименование
	Capacity    int      `json:"Вместимость"` // Имя реквизита в 1С
	RoomType    RoomType `json:"ТипКабинета"` // Имя реквизита в 1С
	Building    Building `json:"Корпус"`      // Имя реквизита в 1С
}

// ClassroomCreatePayload описывает структуру JSON для частичного изменения кабинета (PATCH).
// Использование указателей (*string, *int) критически важно: если фронтенд пришлет только
// {"capacity": 30}, то поля Number, RoomType и Building будут равны nil.
// Это позволяет Go понять, что их обновлять не нужно, а не затирать их пустыми значениями.
type ClassroomUpdatePayload struct {
	Number   *string   `json:"number,omitempty"`    // Необязательный новый номер аудитории
	Capacity *int      `json:"capacity,omitempty"`  // Необязательная новая вместимость
	RoomType *RoomType `json:"room_type,omitempty"` // Необязательный новый тип кабинета
	Building *Building `json:"building,omitempty"`  // Необязательный новый корпус
}

// ODataClassroomUpdate передает измененные реквизиты аудитории в 1С в формате PATCH.
// Теги omitempty указывают Go-серверу полностью исключать непереданные (равные значению по умолчанию)
// поля из результирующего JSON-пакета, чтобы 1С обновила исключительно запрошенные реквизиты.
type ODataClassroomUpdate struct {
	Description string   `json:"Description,omitempty"` // Новое наименование (если менялось)
	Capacity    int      `json:"Вместимость,omitempty"` // Новая вместимость (если менялась)
	RoomType    RoomType `json:"ТипКабинета,omitempty"` // Новый тип кабинета (если менялся)
	Building    Building `json:"Корпус,omitempty"`      // Новый корпус (если менялся)
}
