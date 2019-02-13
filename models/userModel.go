package models

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
	EmailAddress string `json:"emailAddress"`
	UserRole int `json:"userRole"`
}
