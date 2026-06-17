package main

import (
	"log"
	"net/http"

	"academic-booking-api/client"
)

const (
	ODataURL   = "http://localhost:8080/booking/odata/standard.odata/"
	Username   = "administrator"
	Password   = ""
	ServerPort = ":4000"
)

func main() {
	log.Println("Запуск модульной API-обертки...")

	odataClient := client.NewODataClient(ODataURL, Username, Password)
	log.Printf("Клиент OData настроен на адрес: %s\n", odataClient.BaseURL)

	router := setupRoutes(odataClient)

	log.Printf("Go-сервер успешно запущен на http://localhost%s\n", ServerPort)
	if err := http.ListenAndServe(ServerPort, router); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
