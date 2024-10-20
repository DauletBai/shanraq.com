package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}
