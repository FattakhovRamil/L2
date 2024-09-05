package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:

go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type options struct {
	timeout time.Duration
	host    string
	port    string
}

func main() {
	args := os.Args
	ops, err := handlerArgs(args)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		os.Exit(1)
	}
	fmt.Printf("Таймаут подключения: %v\n", ops.timeout)
	fmt.Printf("Хост: %s\n", ops.host)
	fmt.Printf("Порт: %s\n", ops.port)

	// Создаем канал для сигналов
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	err = handlerOps(&ops, signals)
	if err != nil {
		fmt.Println(err)
	}

	// Ожидаем сигнал завершения работы
	<-signals
	fmt.Println("\nПолучен сигнал завершения. Завершение работы...")
}

func handlerArgs(args []string) (options, error) {
	if len(args) < 3 || len(args) > 4 {
		return options{}, fmt.Errorf("неверное количество аргументов: %d", len(args))
	}

	ops := options{
		timeout: 10 * time.Second,
	}

	if len(args) == 4 && strings.HasPrefix(args[1], "--timeout=") {
		parts := strings.Split(args[1], "=")
		if len(parts) != 2 {
			return options{}, fmt.Errorf("некорректный формат строки для таймаута")
		}
		var err error
		ops.timeout, err = time.ParseDuration(parts[1])
		if err != nil {
			return options{}, fmt.Errorf("некорректный формат строки для таймаута: %v", err)
		}
		ops.host = args[2]
		ops.port = args[3]
	} else {
		if len(args) != 3 {
			return options{}, fmt.Errorf("некорректное количество аргументов")
		}
		ops.host = args[1]
		ops.port = args[2]
	}

	return ops, nil
}

func handlerOps(ops *options, signals chan os.Signal) error {
	address := fmt.Sprintf("%s:%s", ops.host, ops.port)
	fmt.Printf("Попытка подключения к %s с таймаутом %v\n", address, ops.timeout)

	timeout := time.After(ops.timeout)
	var conn net.Conn
	var err error

	for {
		select {
		case <-signals:
			return fmt.Errorf("принят сигнал завершения работы")
		case <-timeout:
			return fmt.Errorf("таймаут подключения к %s", address)
		default:
			conn, err = net.Dial("tcp", address)
			if err == nil {
				break
			}
			fmt.Printf("Ошибка подключения: %v. Повторная попытка...\n", err)
			time.Sleep(1 * time.Second)
		}
		if err == nil {
			break
		}
	}

	defer func() {
		conn.Close()
		fmt.Println("Соединение закрыто")
	}()

	// Запускаем горутины для обработки ввода и вывода
	done := make(chan struct{})
	go func() {
		defer close(done)
		userOutputHandler(conn, signals)
	}()

	go func() {
		serverInputHandler(conn, signals)
	}()

	// Ожидаем завершения работы горутин
	<-done
	return nil
}

func userOutputHandler(conn net.Conn, signals chan os.Signal) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-signals:
			return
		default:
			if scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				if _, err := conn.Write([]byte(line + "\n")); err != nil {
					fmt.Println("Ошибка при отправке данных на сервер:", err)
					return
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Ошибка ввода:", err)
				return
			}
		}
	}
}

func serverInputHandler(conn net.Conn, signals chan os.Signal) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-signals:
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Соединение закрыто сервером")
				} else {
					fmt.Println("Ошибка при чтении данных:", err)
				}
				return
			}
			fmt.Print(string(buf[:n]))
		}
	}
}