package models

import (
	"github.com/cvtracker/chaincode/model"
)

type TemplateData struct {
	UserDetails UserDetails
	CurrentPage string
	MessageWarning string
	MessageSuccess string
	CV *model.CVObject
	CVHistory []model.CVHistoryInfo
	Ratings []model.CVRating
	CurrentCVInReview bool
	UserHasCVInReview bool
	CVList map[int] *model.CVObject
}