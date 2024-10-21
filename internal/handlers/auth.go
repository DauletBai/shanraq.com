package handlers

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"shanraq.com/internal/models"
)

var validate = validator.New()

type RegisterInput struct {
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	Birthday  string `json:"birthday" validate:"required,datetime=2006-01-02"`
	FirstName string `json:"first_name" validate:"required,alpha"`
	LastName  string `json:"last_name" validate:"required,alpha"`
	Phone     string `json:"phone" validate:"required,e164"`
	Password  string `json:"password" validate:"required,email"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	// Parsing json
	// ...

	// Data validation
	if err := validate.Struct(input); err != nil {
		// handler
		return
	}

	// Age verification
	birthday, _ := time.Parse("2006-01-02", input.Birthday)
	if age := time.Now().Year() - birthday.Year(); age < 16 {
		// Age err
		return
	}

	// Converting first and last names to the correct case
	input.FirstName = capitalize(input.FirstName)
	input.LastName = capitalize(input.LastName)

	// ...
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return string.ToUpper(string(s[0])) + string.ToLower(s[1:])
}
