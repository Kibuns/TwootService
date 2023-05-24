package Models

import (
	"time"
)

type Twoot struct {
	UserName	string	`json:"username"`
	Content  string    `json:"content"`
	Created time.Time `json:"created"`
}
