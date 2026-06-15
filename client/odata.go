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

// Get выполняет авторизованный Get-запрос к 1С для получения объектов
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

// Post выполняет авторизованный Post-запрос к 1С для создания объекта
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

// Patch выполняет авторизованный PATCH-запрос к 1С для частичного обновления объекта
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
		return fmt.Errorf("ошибка выполнения PATCH-запроса: %v", err)
	}
	defer resp.Body.Close()

	// 1С при успешном обновлении возвращает 200 OK или 204 No Content
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("1С вернула статус %s: %s", resp.Status, string(respBody))
	}

	return nil
}

// Delete выполняет авторизованный DELETE-запрос к 1С для удаления объекта
func (c *ODataClient) Delete(endpoint string) error {
	req, err := http.NewRequest("DELETE", c.BaseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания DELETE-запроса: %v", err)
	}

	req.Header.Set("Authorization", c.AuthHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения DELETE-запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("1С вернула статус %s: %s", resp.Status, string(respBody))
	}

	return nil
}

// Ping проверка доступности OData-сервиса 1С
func (c *ODataClient) Ping() error {
	// GET-запрос к метаданным
	_, err := c.Get("$metadata")
	return err
}
