package handlers

import (
	"net/http"
	"strings"
	"time"
)

var validate = validator.New()

type RegisterInput struct {
	Gender string `json:"gender" validate:"required,oneof=male female"`
	Birthday string `json:"birthday" validate:"required,datetime=2006-01-02"`
	// ...
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
