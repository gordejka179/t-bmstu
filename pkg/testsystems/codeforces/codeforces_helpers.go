package codeforces

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetCodeforcesTaskList(count int) ([]Task, error) {
	result, err := http.Get("https://codeforces.com/problemset")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	doc, err := goquery.NewDocumentFromReader(result.Body)
	if err != nil {
		return nil, err
	}

	var tasks []Task

	doc.Find("table.problems tr").Each(func(i int, s *goquery.Selection) {
		if i < count {
			contestIDStr := s.Find("td").Eq(0).Text()
			problemIDStr := s.Find("td").Eq(1).Text()
			name := s.Find("td").Eq(2).Text()

			contestID := strings.TrimSpace(contestIDStr)
			problemID := strings.TrimSpace(problemIDStr)
			name = strings.TrimSpace(name)

			if contestID != "" && problemID != "" && name != "" {
				idBytes := []byte("cf" + contestID + problemID)
				tasks = append(tasks, Task{
					ID:   base64.StdEncoding.EncodeToString(idBytes),
					Name: name,
				})
			}
		}
	})

	return tasks, nil
}

func SubmitProblem(client *http.Client, submitURL string, fileData url.Values, taskID string) (string, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Добавление пустого файла в multipart форму
	formFile, _ := writer.CreateFormFile("fname", "")
	formFile.Write([]byte(""))

	for key, val := range fileData {
		_ = writer.WriteField(key, val[0])
	}
	writer.Close()

	req, err := http.NewRequest("POST", submitURL, &buffer)
	if err != nil {
		return "0", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "0", err
	}
	defer resp.Body.Close()

	// Формируем URL для получения статуса решения
	statusURL := fmt.Sprintf("https://codeforces.com/problemset/status/%s", taskID)
	statusResp, err := http.Get(statusURL)
	if err != nil {
		fmt.Println("Error:", err)
		return "1", err
	}
	defer statusResp.Body.Close()

	statusDoc, err := goquery.NewDocumentFromReader(statusResp.Body)
	if err != nil {
		log.Fatal(err)
		return "1", err
	}

	// Извлечение статуса решения из таблицы
	table := statusDoc.Find("table.status-frame-datatable")
	if table.Length() > 0 {
		rows := table.Find("tr")
		for i := 1; i < rows.Length(); i++ { // Начинаем с 1, чтобы пропустить заголовок
			row := rows.Eq(i)
			columns := row.Find("td")
			if columns.Length() > 0 {
				status := columns.Eq(3).Text()
				return strings.TrimSpace(status), nil
			}
		}
	}

	return "Unknown", nil
}

func saveToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

type Task struct {
	ID   string
	Name string
}
