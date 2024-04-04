package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./wget <URL>")
		return
	}

	url := os.Args[1]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка при получении страницы: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "ошибка: код состояния %d\n", resp.StatusCode)
		return
	}

	fileName := getFileName(url)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка при создании файла: %v\n", err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка при копировании данных в файл: %v\n", err)
		return
	}

	fmt.Printf("Скачивание завершено. Содержимое сохранено в файл %s\n", fileName)
}

func getFileName(url string) string {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]
	if fileName == "" {
		fileName = "index.html"
	}
	return fileName
}
