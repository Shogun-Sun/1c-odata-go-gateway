package main

import (
	"log"
	"net/http"
	"os"

	"academic-booking-api/internal"
	"academic-booking-api/internal/client"
)

func main() {
	log.Println("Запуск модульной API-обертки...")

	odataURL := os.Getenv("ODATA_URL")
	username := os.Getenv("ODATA_USERNAME")
	password := os.Getenv("ODATA_PASSWORD")
	serverPort := os.Getenv("GO_SERVER_PORT")

	if odataURL == "" {
		log.Fatal("Переменная окружения ODATA_URL не задана!")
	}
	if username == "" {
		log.Fatal("Переменная ODATA_USERNAME не задана!")
	}
	if serverPort == "" {
		serverPort = ":4000"
	}

	odataClient := client.NewODataClient(odataURL, username, password)
	log.Printf("Клиент OData настроен на адрес: %s\n", odataClient.BaseURL)

	router := internal.SetupRoutes(odataClient)

	log.Printf("Go-сервер запущен на порту %s\n", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
