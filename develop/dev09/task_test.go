package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlerURL(t *testing.T) {
	// Создаем тестовый сервер с фиксированным ответом
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer ts.Close()

	// Запускаем тест
	got, err := handlerURL(ts.URL)
	if err != nil {
		t.Fatalf("handlerURL() returned an error: %v", err)
	}

	want := "test content"
	if got != want {
		t.Errorf("handlerURL() = %v; want %v", got, want)
	}
}

func TestConvertContentToPath(t *testing.T) {
	content := "test content"
	url := "http://example.com/page"

	// Временный путь для файла
	tempFilePath := "http_example_com_page.html"
	defer os.Remove(tempFilePath)

	err := convertContentToPath(content, url)
	if err != nil {
		t.Fatalf("convertContentToPath() returned an error: %v", err)
	}

	// Проверка существования файла
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		t.Errorf("File %v does not exist", tempFilePath)
	}
}

func TestConvertURLToFilePath(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"http://example.com", "http_example_com.html"},
		{"https://example.com/path/to/resource", "https_example_com_path_to_resource.html"},
		{"ftp://example.com/file.txt", "ftp_example_com_file_txt.html"},
	}

	for _, test := range tests {
		got, err := convertURLToFilePath(test.url)
		if err != nil {
			t.Fatalf("convertURLToFilePath(%q) returned an error: %v", test.url, err)
		}

		if got != test.expected {
			t.Errorf("convertURLToFilePath(%q) = %q; want %q", test.url, got, test.expected)
		}
	}
}
