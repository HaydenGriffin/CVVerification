package models

import (
	"github.com/cvverification/chaincode/model"
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
	Review model.CVReview
	Reviews []model.CVReview
	CurrentCVInReview bool
	CurrentCVHash string
	UserHasCVInReview bool
	CVList map[int] *model.CVObject
}