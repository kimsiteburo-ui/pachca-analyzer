package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ApiResponse reflects the top-level structure of the API response
type ApiResponse struct {
	Data []User `json:"data"`
}

type User struct {
	ID          int     `json:"id"`
	Email       string  `json:"email"`
	Nickname    string  `json:"nickname"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Role        string  `json:"role"`
	PhoneNumber *string `json:"phone_number"`
	Department  *string `json:"department"`
	Title       *string `json:"title"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiToken := os.Getenv("PACHCA_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Переменная окружения PACHCA_API_TOKEN не установлена")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.pachca.com/api/shared/v1/users", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка API: %s\n", string(body))
		return
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	for _, user := range apiResponse.Data {
		phoneNumber := ""
		if user.PhoneNumber != nil {
			phoneNumber = *user.PhoneNumber
		}
		department := ""
		if user.Department != nil {
			department = *user.Department
		}
		title := ""
		if user.Title != nil {
			title = *user.Title
		}
		fmt.Printf("ID: %d, Email: %s, Nickname: %s, FullName: %s %s, Role: %s, Phone: %s, Department: %s, Title: %s\n", user.ID, user.Email, user.Nickname, user.FirstName, user.LastName, user.Role, phoneNumber, department, title)
	}
}
