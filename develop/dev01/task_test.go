package main

import (
	"testing"
)

// TestMainSuccess проверяет корректное выполнение программы без ошибок
func TestPrintCurrentTimeSuccess(t *testing.T) {
	// Сохраняем исходную функцию ntp.Time и восстанавливаем её после теста
	host := "0.beevik-ntp.pool.ntp.org"
	_, err := printCurrentTime(host)
	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

}

func TestPrintCurrentTimeError(t *testing.T) {
	// Сохраняем исходную функцию ntp.Time и восстанавливаем её после теста
	host := "0.beevik-ntp.pool.ntp.org1"
	result, err := printCurrentTime(host)
	if err == nil {
		t.Fatalf("ожидалось ошибочное выполнение, но получен результат: %v", result)
	}

}

// TestMainError проверяет корректное выполнение программы при ошибке
