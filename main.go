package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
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

func mustSetCellValue(f *excelize.File, sheet, axis string, value interface{}) {
	if err := f.SetCellValue(sheet, axis, value); err != nil {
		panic(err)
	}
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

	f := excelize.NewFile()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set value of a cell.
	mustSetCellValue(f, "Sheet1", "A1", "ID")
	mustSetCellValue(f, "Sheet1", "B1", "Email")
	mustSetCellValue(f, "Sheet1", "C1", "Nickname")
	mustSetCellValue(f, "Sheet1", "D1", "FirstName")
	mustSetCellValue(f, "Sheet1", "E1", "LastName")
	mustSetCellValue(f, "Sheet1", "F1", "Role")
	mustSetCellValue(f, "Sheet1", "G1", "PhoneNumber")
	mustSetCellValue(f, "Sheet1", "H1", "Department")
	mustSetCellValue(f, "Sheet1", "I1", "Title")

	for i, user := range apiResponse.Data {
		row := i + 2
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

		mustSetCellValue(f, "Sheet1", fmt.Sprintf("A%d", row), user.ID)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("B%d", row), user.Email)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("C%d", row), user.Nickname)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("D%d", row), user.FirstName)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("E%d", row), user.LastName)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("F%d", row), user.Role)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("G%d", row), phoneNumber)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("H%d", row), department)
		mustSetCellValue(f, "Sheet1", fmt.Sprintf("I%d", row), title)
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs("pachca_users.xlsx"); err != nil {
		fmt.Println(err)
	}
}
