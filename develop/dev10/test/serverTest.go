package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Обработка соединений с клиентами
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Клиент подключен:", conn.RemoteAddr())

	// Создаем новый сканер для чтения данных от клиента
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fmt.Println("Получено от клиента:", line)

		// Отправляем обратно клиенту те же данные
		_, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("Ошибка отправки данных:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения данных:", err)
	}
}

func main() {
	// Параметры сервера
	port := "8080"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту", port)

	// Основной цикл сервера
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка принятия соединения:", err)
			continue
		}

		go handleConnection(conn)
	}
}
