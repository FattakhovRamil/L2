package main

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов

Методы API:
	POST /create_event
	POST /update_event
	POST /delete_event

	GET /events_for_day
	GET /events_for_week
	GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Event представляет одно событие.
type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}

// Storage для хранения событий в памяти.
var events = []Event{}
var nextID = 1

// serializeJSON сериализует объект в JSON.
func serializeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// deserializeJSON десериализует JSON в объект.
func deserializeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// parseEventParams парсирует параметры для создания или обновления события.
func parseEventParams(r *http.Request) (Event, error) {
	var event Event
	err := r.ParseForm()
	if err != nil {
		return event, err
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return event, errors.New("invalid user_id")
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return event, errors.New("invalid date format")
	}

	event.UserID = userID
	event.Title = r.FormValue("title")
	event.Date = date

	return event, nil
}

// parseQueryParams парсирует параметры для GET-запросов.
func parseQueryParams(r *http.Request) (time.Time, error) {
	// userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	// if err != nil {
	// 	return time.Time{}, errors.New("invalid user_id")
	// }

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}

	return date, nil
}

// createEventHandler обрабатывает создание события.
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parseEventParams(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	event.ID = nextID
	nextID++
	events = append(events, event)

	response, _ := serializeJSON(map[string]string{"result": "Event created"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// updateEventHandler обрабатывает обновление события.
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parseEventParams(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	for i, e := range events {
		if e.ID == event.ID {
			events[i] = event
			response, _ := serializeJSON(map[string]string{"result": "Event updated"})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
			return
		}
	}

	http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
}

// deleteEventHandler обрабатывает удаление события.
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, `{"error": "invalid event ID"}`, http.StatusBadRequest)
		return
	}

	for i, e := range events {
		if e.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			response, _ := serializeJSON(map[string]string{"result": "Event deleted"})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
			return
		}
	}

	http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
}

// eventsForDayHandler обрабатывает запрос событий за день.
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	var result []Event
	for _, e := range events {
		if e.Date.Year() == date.Year() && e.Date.YearDay() == date.YearDay() {
			result = append(result, e)
		}
	}

	response, _ := serializeJSON(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// eventsForWeekHandler обрабатывает запрос событий за неделю.
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	startOfWeek := date.Add(-time.Duration(date.Weekday()) * 24 * time.Hour)
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	var result []Event
	for _, e := range events {
		if !e.Date.Before(startOfWeek) && e.Date.Before(endOfWeek) {
			result = append(result, e)
		}
	}

	response, _ := serializeJSON(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// eventsForMonthHandler обрабатывает запрос событий за месяц.
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Определяем начало и конец месяца
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 3, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var result []Event
	for _, e := range events {
		// Проверяем, попадает ли событие в заданный месяц
		if !e.Date.Before(startOfMonth) && e.Date.Before(endOfMonth) {
			result = append(result, e)
		}
	}

	response, _ := serializeJSON(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// loggingMiddleware логирует запросы.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, duration)
	})
}

func main() {
	http.Handle("/create_event", loggingMiddleware(http.HandlerFunc(createEventHandler)))
	http.Handle("/update_event", loggingMiddleware(http.HandlerFunc(updateEventHandler)))
	http.Handle("/delete_event", loggingMiddleware(http.HandlerFunc(deleteEventHandler)))
	http.Handle("/events_for_day", loggingMiddleware(http.HandlerFunc(eventsForDayHandler)))
	http.Handle("/events_for_week", loggingMiddleware(http.HandlerFunc(eventsForWeekHandler)))
	http.Handle("/events_for_month", loggingMiddleware(http.HandlerFunc(eventsForMonthHandler)))

	port := ":8080" // Порт сервера
	fmt.Printf("Server listening on port%s\n", port)
	http.ListenAndServe(port, nil)
}
