package models

type User struct {
	ID       int    `json:"-" db:"id"`
	About    string `json:"about" db:"about"`
	Email    string `json:"email" db:"email"`
	Fullname string `json:"fullname" db:"fullname"`
	Nickname string `json:"nickname" db:"nickname"`
}

// func (user *User) ToLowerUser() {
// 	user.Email = strings.ToLower(user.Email)
// 	user.Fullname = strings.ToLower(user.Fullname)
// 	user.Nickname = strings.ToLower(user.Nickname)
// }

// func (user *User) IsTheSameAs(another User) bool {
// 	return strings.ToLower(user.Email) == strings.ToLower(another.Email) &&
// 		strings.ToLower(user.Nickname) == strings.ToLower(another.Nickname)
// }
