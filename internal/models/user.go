package models

import (
	"fmt"
	"strings"
)

// User model
type User struct {
	ID       int    `json:"-" db:"id"`
	About    string `json:"about" db:"about"`
	Email    string `json:"email" db:"email"`
	Fullname string `json:"fullname" db:"fullname"`
	Nickname string `json:"nickname" db:"nickname"`
}

// FillEmpty fill empty
func (user *User) FillEmpty(another User) {
	if strings.Trim(user.Email, " ") == "" {
		user.Email = another.Email
	}
	if strings.Trim(user.Fullname, " ") == "" {
		user.Fullname = another.Fullname
	}
	if strings.Trim(user.About, " ") == "" {
		user.About = another.About
	}
}

// Print for debug
func (user *User) Print() {
	fmt.Println("-------User-------")
	fmt.Println("--About:", user.About)
	fmt.Println("--Email:", user.Email)
	fmt.Println("--Fullname:", user.Fullname)
	fmt.Println("--Nickname:", user.Nickname)
}
