package models

type TemplateData struct {
	TxId string
	CurrentUser User
	CurrentPage string
	LoggedInFlag bool
	MessageWarning string
	MessageSuccess string
}