package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	ODataURL = "http://localhost:8080/odata_base/odata/standard.odata/"
	Username = "administrator"
	Password = ""

	ServerPort = ":4000"
)

type ODataClient struct {
	BaseURL    string
	AuthHeader string
	Client     *http.Client
}

func NewODataClient(baseURL, user, pass string) *ODataClient {
	authStr := fmt.Sprintf("%s:%s", user, pass)

	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))

	return &ODataClient{
		BaseURL:    baseURL,
		AuthHeader: "Basic " + encodedAuth,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func main() {
	fmt.Println("Инициализация API")

	odataClient := NewODataClient(ODataURL, Username, Password)
	log.Printf("Клиент OData настроен на: %s\n", odataClient.BaseURL)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	log.Printf("Go-сервер успешно запущен на http://localhost%s\n", ServerPort)
	err := http.ListenAndServe(ServerPort, mux)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
