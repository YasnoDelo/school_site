package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/YasnoDelo/school_site/internal/server"
)

func TestHandleWelcome(t *testing.T) {
	// Создаем новый экземпляр сервера
	testServer := server.NewServer()

	// Создаем новый HTTP-запрос
	request, err := http.NewRequest("GET", "/hi", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый HTTP-рекордер (для записи ответа)
	recorder := httptest.NewRecorder()

	// Вызываем метод ServeHTTP для обработки запроса
	testServer.ServeHTTP(recorder, request)

	// Проверяем, что код ответа равен 200
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем, что ответ содержит приветственное сообщение
	expected := "Welcome to the Best HTTP Server Ever!"
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestServeImageAtRoot(t *testing.T) {
	// Создаем новый экземпляр сервера
	testServer := server.NewServer()

	// Создаем новый HTTP-запрос
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый HTTP-рекордер (для записи ответа)
	recorder := httptest.NewRecorder()

	// Вызываем метод ServeHTTP для обработки запроса
	testServer.ServeHTTP(recorder, request)

	// Проверяем, что код ответа равен 200
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что ответ содержит данные изображения (должен быть двоичный контент)
	contentType := recorder.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		t.Errorf("handler returned wrong content type: got %v want image/*", contentType)
	}

	// Проверяем, что тело ответа не пустое
	body, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Errorf("handler returned empty body for image")
	}
}

func TestHandleMetricUpdate(t *testing.T) {
	// Создаем новый экземпляр сервера
	testServer := server.NewServer()

	// Создаем новый HTTP-запрос для обновления метрики
	request, err := http.NewRequest("POST", "/update/gauge/HeapAlloc/123.45", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый HTTP-рекордер (для записи ответа)
	recorder := httptest.NewRecorder()

	// Вызываем метод ServeHTTP для обработки запроса
	testServer.ServeHTTP(recorder, request)

	// Проверяем, что код ответа равен 200
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем, что метрика была правильно обновлена
	metrics := testServer.GetMetrics()
	metric, exists := metrics["HeapAlloc"]
	if !exists {
		t.Errorf("metric not found: %v", "HeapAlloc")
	} else {
		if metric.Type != "gauge" {
			t.Errorf("metric type mismatch: got %v want %v", metric.Type, "gauge")
		}
		if metric.Value != 123.45 {
			t.Errorf("metric value mismatch: got %v want %v", metric.Value, 123.45)
		}
	}
}

func TestInvalidMetricUpdate(t *testing.T) {
	// Создаем новый экземпляр сервера
	testServer := server.NewServer()

	// Создаем новый HTTP-запрос с некорректным типом метрики
	request, err := http.NewRequest("POST", "/update/unknown/HeapAlloc/123.45", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый HTTP-рекордер (для записи ответа)
	recorder := httptest.NewRecorder()

	// Вызываем метод ServeHTTP для обработки запроса
	testServer.ServeHTTP(recorder, request)

	// Проверяем, что код ответа равен 400 (Bad Request)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestNotFound(t *testing.T) {
	// Создаем новый экземпляр сервера
	testServer := server.NewServer()

	// Создаем новый HTTP-запрос с несуществующим маршрутом
	request, err := http.NewRequest("GET", "/unknown", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый HTTP-рекордер (для записи ответа)
	recorder := httptest.NewRecorder()

	// Вызываем метод ServeHTTP для обработки запроса
	testServer.ServeHTTP(recorder, request)

	// Проверяем, что код ответа равен 404 (Not Found)
	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
