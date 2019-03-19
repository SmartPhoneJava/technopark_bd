package models

import "time"

type Thread struct {
	Author  string    `json:"author" db:"author"`
	Created time.Time `json:"created" db:"created"`
	Forum   string    `json:"forum" db:"forum"`
	ID      int       `json:"id" db:"id"`
	Message string    `json:"message" db:"message"`
	Slug    string    `json:"slug" db:"slug"`
	Title   string    `json:"title" db:"title"`
	Votes   int       `json:"-" db:"votes"`
}
