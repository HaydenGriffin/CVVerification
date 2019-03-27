package models

import (
	"github.com/cvtracker/chaincode/model"
)

type TemplateData struct {
	UserDetails UserDetails
	CurrentPage string
	MessageWarning string
	MessageSuccess string
	CVInfo CVDisplayInfo
}

type CVDisplayInfo struct {
	CV *model.CVObject
	CVHistory []model.CVHistoryInfo
	Rating model.CVRating
	Ratings []model.CVRating
	CurrentCVInReview bool
	CurrentCVHash string
	UserHasCVInReview bool
	CVList map[int] *model.CVObject
}