package models

import "github.com/cvtracker/service"

type TemplateData struct {
	TxId string
	CurrentUser User
	CurrentPage string
	LoggedInFlag bool
	MessageWarning string
	MessageSuccess string
	CV service.CVObject
	IsCVRatable bool
	CVList map[int] service.CVObject
}