package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCreateEventHandler проверяет обработчик создания события.
func setup() {
	events = []Event{}
	nextID = 1
}

// TestCreateEventHandler проверяет обработчик создания события.
func TestCreateEventHandler(t *testing.T) {
	setup() // Очищаем глобальные переменные.

	// Создаем новый запрос.
	req, err := http.NewRequest("POST", "/create_event", bytes.NewBufferString("user_id=1&date=2024-07-20&title=New%20Event"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Создаем тестовый сервер.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)

	// Выполняем запрос.
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем результат.
	var result map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}
	if result["result"] != "Event created" {
		t.Errorf("handler returned unexpected body: got %v want %v", result["result"], "Event created")
	}
}

// TestEventsForDayHandler проверяет обработчик запроса событий за день.
func TestEventsForDayHandler(t *testing.T) {
	setup() // Очищаем глобальные переменные.

	// Сначала создаем событие для проверки.
	TestCreateEventHandler(t)

	// Создаем запрос на получение событий за день.
	req, err := http.NewRequest("GET", "/events_for_day?date=2024-07-20", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый сервер.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventsForDayHandler)

	// Выполняем запрос.
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем результат.
	var result []Event
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}

	// Проверяем, что результат содержит правильные события за день.
	expectedEvent := Event{
		ID:     1, // ID должен быть таким же, как в тесте создания события
		UserID: 1,
		Title:  "New Event",
		Date:   time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
	}

	if len(result) != 1 {
		t.Errorf("handler returned wrong number of events: got %v want 1", len(result))
		return
	}

	if result[0] != expectedEvent {
		t.Errorf("handler returned unexpected body: got %+v want %+v", result[0], expectedEvent)
	}
}
