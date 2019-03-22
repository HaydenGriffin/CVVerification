package models

import (
	"github.com/cvtracker/chaincode/model"
)

type TemplateData struct {
	TxId string
	UserDetails UserDetails
	CurrentPage string
	IsAdmin bool
	IsApplicant bool
	MessageWarning string
	MessageSuccess string
	CV *model.CVObject
	Ratings []model.CVRating
	IsCVRatable bool
	CVList map[int] model.CVObject
}