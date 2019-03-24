package models

import (
	"fmt"
	"time"
)

type Post struct {
	Author   string    `json:"author" db:"author"`
	Created  time.Time `json:"created" db:"created"`
	Forum    string    `json:"forum" db:"forum"`
	ID       int       `json:"id" db:"id"`
	IsEdited bool      `json:"-" db:"isEdited"`
	Message  string    `json:"message" db:"message"`
	Parent   int       `json:"parent" db:"parent"`
	Thread   int       `json:"thread" db:"thread"`
	Path     string    `json:"path" db:"path"`
}

func (post *Post) Print() {
	fmt.Println("-------Post-------")
	fmt.Println("--ID:", post.ID)
	fmt.Println("--Parent:", post.Parent)
	fmt.Println("--Path:", post.Path)
	fmt.Println("--Created:", post.Created)
}
