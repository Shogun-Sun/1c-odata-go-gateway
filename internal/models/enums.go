package models

// RoomType определяет тип назначения учебной аудитории.
type RoomType string

const (
	RoomComputer RoomType = "КомпьютерныйКласс"
	RoomLecture  RoomType = "Лекционный"
	RoomLab      RoomType = "Лаборатория"
)

// IsValid возвращает true, если назначение кабинета поддерживается системой.
func (rt RoomType) IsValid() bool {
	switch rt {
	case RoomComputer, RoomLecture, RoomLab:
		return true
	}
	return false
}

// Building определяет учебный корпус или площадку расположения кабинета.
type Building string

const (
	BuildingMain   Building = "ПерваяПлощадка"
	BuildingSecond Building = "ВтораяПлощадка"
	BuildingSormov Building = "Общежитие"
)

// IsValid проверяет принадлежность здания к списку официальных корпусов колледжа.
func (b Building) IsValid() bool {
	switch b {
	case BuildingMain, BuildingSecond, BuildingSormov:
		return true
	}
	return false
}

// Position определяет должность преподавателя в виде строго типизированной строки.
type Position string

const (
	PositionTeacher Position = "Преподаватель"
)

// IsValid сопоставляет строку со списком разрешенных штатных должностей 1С.
func (p Position) IsValid() bool {
	switch p {
	case PositionTeacher:
		return true
	}
	return false
}

// LessonType определяет вид занятия
type LessonType string

const (
	LessonLecture      LessonType = "Лекция"
	LessonLaboratory   LessonType = "Лабораторная"
	LessonExam         LessonType = "Экзамен"
	LessonConsultation LessonType = "Консультация"
)

// IsValid проверяет, входит ли переданный вид занятия в список допустимых в 1С.
func (lt LessonType) IsValid() bool {
	switch lt {
	case LessonLecture, LessonLaboratory, LessonExam, LessonConsultation:
		return true
	}

	return false
}
