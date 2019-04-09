package model

type UserDetails struct {
	Id int `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	EmailAddress string `json:"emailAddress"`
	UploadedCV bool `json:"uploadedCV"`
}