package models

type TemplateData struct {
	CurrentUser User
	CurrentPage string
	LoggedInFlag bool
	MessageWarning string
	MessageSuccess string
}