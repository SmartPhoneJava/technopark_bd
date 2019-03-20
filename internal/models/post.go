package models

import "time"

type Post struct {
	Author   string    `json:"author" db:"author"`
	Created  time.Time `json:"created" db:"created"`
	Forum    string    `json:"forum" db:"forum"`
	ID       int       `json:"id" db:"id"`
	IsEdited bool      `json:"-" db:"isEdited"`
	Message  string    `json:"message" db:"message"`
	Parent   int       `json:"parent" db:"parent"`
	Thread   int       `json:"thread" db:"thread"`
	Path     string    `json:"-" db:"path"`
}
