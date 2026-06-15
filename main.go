package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/handlers"
)

const (
	ODataURL   = "http://localhost:8080/odata_base/odata/standard.odata/"
	Username   = "administrator"
	Password   = ""
	ServerPort = ":4000"
)

func main() {
	log.Println("Запуск модульной API-обертки...")

	// Инициализируем клиент
	odataClient := client.NewODataClient(ODataURL, Username, Password)
	log.Printf("Клиент OData настроен на адрес: %s\n", odataClient.BaseURL)

	// Инициализация хэндлеров
	groupHandler := handlers.NewGroupHandler(odataClient)
	classroomHandler := handlers.NewClassroomHandler(odataClient)
	teacherHandler := handlers.NewTeacherHandler(odataClient)
	departmentHandler := handlers.NewDepartmentHandler(odataClient)
	disciplineHandler := handlers.NewDisciplineHandler(odataClient)

	// Подключение роутера для эндпоинтов
	mux := http.NewServeMux()

	// Эндпоинт проверки работоспособности
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		// Проверка работаспособности API
		apiStatus := "OK"

		// Проверка доступности 1С
		err := odataClient.Ping()
		oneCStatus := "OK"
		if err != nil {
			oneCStatus = fmt.Sprintf("FAILED (%v)", err)
		}

		// Сборка ответа
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable) // 503, если 1С недоступна
		}

		response := map[string]string{
			"api": apiStatus,
			"1c":  oneCStatus,
		}

		json.NewEncoder(w).Encode(response)
	})

	// Регистрация эндпоинтов для групп
	mux.HandleFunc("GET /api/v1/groups", groupHandler.GetGroups)
	mux.HandleFunc("POST /api/v1/groups", groupHandler.CreateGroup)
	mux.HandleFunc("PATCH /api/v1/groups/{id}", groupHandler.UpdateGroup)
	mux.HandleFunc("DELETE /api/v1/groups/{id}", groupHandler.DeleteGroup)

	// Регистрация эндпоинтов кабинетов
	mux.HandleFunc("GET /api/v1/classrooms", classroomHandler.GetClassrooms)
	mux.HandleFunc("POST /api/v1/classrooms", classroomHandler.CreateClassroom)
	mux.HandleFunc("PATCH /api/v1/classrooms/{id}", classroomHandler.UpdateClassroom)
	mux.HandleFunc("DELETE /api/v1/classrooms/{id}", classroomHandler.DeleteClassroom)

	// Регистрация эндпоинтов преподавателей
	mux.HandleFunc("GET /api/v1/teachers", teacherHandler.GetTeachers)
	mux.HandleFunc("POST /api/v1/teachers", teacherHandler.CreateTeacher)
	mux.HandleFunc("PATCH /api/v1/teachers/{id}", teacherHandler.UpdateTeacher)
	mux.HandleFunc("DELETE /api/v1/teachers/{id}", teacherHandler.DeleteTeacher)

	// Регистрация эндпоинтов кафедры
	mux.HandleFunc("GET /api/v1/departments", departmentHandler.GetDepartments)
	mux.HandleFunc("POST /api/v1/departments", departmentHandler.CreateDepartment)
	mux.HandleFunc("PATCH /api/v1/departments/{id}", departmentHandler.UpdateDepartment)
	mux.HandleFunc("DELETE /api/v1/departments/{id}", departmentHandler.DeleteDepartment)

	// Регистрация эндпоинтов дисциплин
	mux.HandleFunc("GET /api/v1/disciplines", disciplineHandler.GetDisciplines)
	mux.HandleFunc("POST /api/v1/disciplines", disciplineHandler.CreateDiscipline)
	mux.HandleFunc("PATCH /api/v1/disciplines/{id}", disciplineHandler.UpdateDiscipline)
	mux.HandleFunc("DELETE /api/v1/disciplines/{id}", disciplineHandler.DeleteDiscipline)

	// Запуск сервера
	log.Printf("Go-сервер успешно запущен на http://localhost%s\n", ServerPort)
	err := http.ListenAndServe(ServerPort, mux)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
