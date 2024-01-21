package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const secretKey = "aB3dE4gH"

type ReviewRequest struct {
	Id     int `json:"Id"`
	UserID int `json:"id_user"`
}

type ReviewResult struct {
	Result bool `json:"result"`
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/asyncProcess", handleReview)
	log.Println("Сервер был запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ReviewRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Message: "Запрос не выполнен",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "Запрос выполнен",
	}
	json.NewEncoder(w).Encode(response)

	go processReview(req.Id, req.UserID)
}

func processReview(Id int, id_user int) {
	delay := rand.Intn(6) + 5
	log.Printf("Request ID %d выполняется с задержкой %d секунд", Id, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	rand.Seed(time.Now().UnixNano())

	result := rand.Intn(2) == 1
	log.Println(result)
	sendResult(Id, id_user, result)
}

func sendResult(Id int, id_user int, result bool) error {
	reviewResult := struct {
		Result    bool   `json:"result"`
		UserID    int    `json:"id_user"`
		SecretKey string `json:"secretKey"`
	}{
		Result:    result,
		UserID:    id_user,
		SecretKey: secretKey,
	}

	jsonData, err := json.Marshal(reviewResult)
	if err != nil {
		return fmt.Errorf("Ошибка при маршалинге JSON данных: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://127.0.0.1:8000/update/%d/", Id), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Ошибка при создании PUT-запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Ошибка при выполнении PUT-запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка при отправке результата: код состояния %d", resp.StatusCode)
	}

	log.Printf("Отправлено Request ID [%d] - Результат: %t", Id, result)

	return nil
}
