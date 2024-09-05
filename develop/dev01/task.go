package main

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

var ntpTime = ntp.Time

// printCurrentTime печатает текущее время
func printCurrentTime(host string) (time.Time, error) {
	t, err := ntpTime(host)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// printError печатает ошибку в STDERR

func printError(err error) {

	fmt.Fprintln(os.Stderr, "Ошибка:", err)

}

func main() {

	host := "0.beevik-ntp.pool.ntp.org"
	result, err := printCurrentTime(host)

	fmt.Println("Текущее время: ", result)

	if err != nil {
		printError(err)
		os.Exit(1)
	}
	os.Exit(0)
}
