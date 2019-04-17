package model

import (
	"github.com/cvverification/chaincode/model"
)

type Data struct {
	UserDetails UserDetails
	AccountType string
	CurrentPage string
	MessageWarning string
	MessageSuccess string
	CVInfo CVDisplayInfo
	PrivateKey string
}

type CVDisplayInfo struct {
	CurrentCVID string
	CV *model.CVObject
	CVHistory []CVHistoryInfo
	ReviewInfo CVReviewInfo
	CVList map[string]model.CVObject
	ApplicantDetails UserDetails
}

type CVReviewInfo struct{
	Status string
	Reviews []model.CVReview
}

type CVHistoryInfo struct {
	Index int
	CVID string
	CV *model.CVObject
}
