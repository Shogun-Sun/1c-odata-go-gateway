package main

import (
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

	// 1. Инициализируем клиент из пакета client
	odataClient := client.NewODataClient(ODataURL, Username, Password)
	log.Printf("Клиент OData настроен на адрес: %s\n", odataClient.BaseURL)

	// 2. Инициализируем хендлер групп и передаем ему клиент
	groupHandler := handlers.NewGroupHandler(odataClient)

	// 3. Роутер
	mux := http.NewServeMux()

	// ping
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	// 4. РЕГИСТРИРУЕМ ЭНДПОИНТ ДЛЯ ГРУПП
	mux.HandleFunc("GET /api/v1/groups", groupHandler.GetGroups)
	mux.HandleFunc("POST /api/v1/groups", groupHandler.CreateGroup)

	// 5. Старт сервера
	log.Printf("Go-сервер успешно запущен на http://localhost%s\n", ServerPort)
	err := http.ListenAndServe(ServerPort, mux)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
