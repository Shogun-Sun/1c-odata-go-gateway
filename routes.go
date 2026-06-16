package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"academic-booking-api/client"
	"academic-booking-api/handlers"
)

// setupRoutes инициализирует хендлеры и регистрирует все эндпоинты API,
// возвращая настроенный http.Handler
func setupRoutes(odataClient *client.ODataClient) http.Handler {
	mux := http.NewServeMux()

	groupHandler := handlers.NewGroupHandler(odataClient)
	classroomHandler := handlers.NewClassroomHandler(odataClient)
	departmentHandler := handlers.NewDepartmentHandler(odataClient)
	disciplineHandler := handlers.NewDisciplineHandler(odataClient)
	teacherHandler := handlers.NewTeacherHandler(odataClient)

	// Эндпоинт проверки работоспособности (Health Check)
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		apiStatus := "OK"
		oneCStatus := "OK"

		if err := odataClient.Ping(); err != nil {
			oneCStatus = fmt.Sprintf("FAILED (%v)", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if oneCStatus != "OK" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		response := map[string]string{
			"api": apiStatus,
			"1c":  oneCStatus,
		}

		json.NewEncoder(w).Encode(response)
	})

	// Учебные группы
	mux.HandleFunc("GET /api/v1/groups", groupHandler.GetGroups)
	mux.HandleFunc("POST /api/v1/groups", groupHandler.CreateGroup)
	mux.HandleFunc("PATCH /api/v1/groups/{id}", groupHandler.UpdateGroup)
	mux.HandleFunc("DELETE /api/v1/groups/{id}", groupHandler.DeleteGroup)

	// Кабинеты
	mux.HandleFunc("GET /api/v1/classrooms", classroomHandler.GetClassrooms)
	mux.HandleFunc("POST /api/v1/classrooms", classroomHandler.CreateClassroom)
	mux.HandleFunc("PATCH /api/v1/classrooms/{id}", classroomHandler.UpdateClassroom)
	mux.HandleFunc("DELETE /api/v1/classrooms/{id}", classroomHandler.DeleteClassroom)

	// Преподаватели
	mux.HandleFunc("GET /api/v1/teachers", teacherHandler.GetTeachers)
	mux.HandleFunc("POST /api/v1/teachers", teacherHandler.CreateTeacher)
	mux.HandleFunc("PATCH /api/v1/teachers/{id}", teacherHandler.UpdateTeacher)
	mux.HandleFunc("DELETE /api/v1/teachers/{id}", teacherHandler.DeleteTeacher)

	// Кафедры
	mux.HandleFunc("GET /api/v1/departments", departmentHandler.GetDepartments)
	mux.HandleFunc("POST /api/v1/departments", departmentHandler.CreateDepartment)
	mux.HandleFunc("PATCH /api/v1/departments/{id}", departmentHandler.UpdateDepartment)
	mux.HandleFunc("DELETE /api/v1/departments/{id}", departmentHandler.DeleteDepartment)

	// Дисциплины
	mux.HandleFunc("GET /api/v1/disciplines", disciplineHandler.GetDisciplines)
	mux.HandleFunc("POST /api/v1/disciplines", disciplineHandler.CreateDiscipline)
	mux.HandleFunc("PATCH /api/v1/disciplines/{id}", disciplineHandler.UpdateDiscipline)
	mux.HandleFunc("DELETE /api/v1/disciplines/{id}", disciplineHandler.DeleteDiscipline)

	return mux
}
