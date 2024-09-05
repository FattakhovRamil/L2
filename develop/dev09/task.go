package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Некорректное количество аргументов")
		return
	}
	url := args[1]
	content, err := handlerURL(url)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		os.Exit(1)
	}
	convertContentToPath(content, url)
}

func handlerURL(url string) (string, error) {
	respons, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer func() {
		respons.Body.Close()
	}()

	body, err := io.ReadAll(respons.Body)

	if err != nil {
		//fmt.Println("Ошибка чтения:", err)
		return "", err
	}

	//fmt.Println("Статус запроса:", respons.Status)
	//fmt.Println("Body:")
	return string(body), nil
}

func convertContentToPath(content, pathToSave string) error {
	filePath, err := convertURLToFilePath(pathToSave)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	defer func() {
		file.Close()
	}()

	if err != nil {
		return err
	}

	_, err = file.WriteString(content)

	if err != nil {
		return err
	}

	//fmt.Println("Содержимое успешно записано в файл", filePath)
	return nil
}

func convertURLToFilePath(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("пустой путь к файлу")
	}

	filePath = strings.ReplaceAll(filePath, "://", "_")
	filePath = strings.ReplaceAll(filePath, ".", "_")
	filePath = strings.ReplaceAll(filePath, "/", "_")
	filePath = path.Clean(filePath)

	filePath += ".html"

	return filePath, nil
}
