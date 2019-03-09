package models

type UserDetails struct {
	Id int `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	EmailAddress string `json:"emailAddress"`
	ProfileHash string `json:"profileHash"`
}