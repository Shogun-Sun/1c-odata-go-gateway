// Package client реализует низкоуровневый HTTP-клиент для работы с интерфейсом OData 1С.
// Обеспечивает автоматическую сборку заголовков авторизации, управление таймаутами
// и унифицированную обработку статус-кодов ответов платформы 1С:Предприятие.
package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ODataClient инкапсулирует параметры подключения и HTTP-клиент для отправки запросов в 1С.
type ODataClient struct {
	BaseURL    string       // Корневой URL OData-интерфейса (включая имя информационной базы)
	AuthHeader string       // Предварительно собранный заголовок Basic-авторизации
	Client     *http.Client // Настроенный HTTP-клиент с управлением таймаутами
}

// NewODataClient инициализирует и возвращает новый экземпляр ODataClient.
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

// Get выполняет авторизованный GET-запрос к 1С для получения сырых байтовых данных.
func (c *ODataClient) Get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", c.AuthHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к 1С: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("1С вернула статус: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// Post выполняет авторизованный POST-запрос к 1С для создания нового объекта.
func (c *ODataClient) Post(endpoint string, body interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации JSON: %v", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания POST-запроса: %v", err)
	}

	req.Header.Set("Authorization", c.AuthHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения POST-запроса к 1С: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("1С вернула статус %s: %s", resp.Status, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

// Patch выполняет авторизованный PATCH-запрос к 1С для частичного обновления существующего объекта.
func (c *ODataClient) Patch(endpoint string, body interface{}) error {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %v", err)
	}

	req, err := http.NewRequest("PATCH", c.BaseURL+endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("ошибка создания PATCH-запроса: %v", err)
	}

	req.Header.Set("Authorization", c.AuthHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения PATCH-запроса к 1С: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("1С вернула статус %s: %s", resp.Status, string(respBody))
	}

	return nil
}

// Delete выполняет авторизованный DELETE-запрос к 1С для удаления объекта по его идентификатору.
func (c *ODataClient) Delete(endpoint string) error {
	req, err := http.NewRequest("DELETE", c.BaseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания DELETE-запроса: %v", err)
	}

	req.Header.Set("Authorization", c.AuthHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения DELETE-запроса к 1С: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("1С вернула статус %s: %s", resp.Status, string(respBody))
	}

	return nil
}

// Ping проверяет доступность OData-сервиса 1С путем запроса метаданных системы.
func (c *ODataClient) Ping() error {
	_, err := c.Get("$metadata")
	return err
}
