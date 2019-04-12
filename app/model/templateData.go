package model

import (
	"github.com/cvverification/chaincode/model"
)

type Data struct {
	UserDetails UserDetails
	CurrentPage string
	MessageWarning string
	MessageSuccess string
	CVInfo CVDisplayInfo
}

type CVDisplayInfo struct {
	CurrentCVID string
	CV *model.CVObject
	CVHistory []CVHistoryInfo
	VerifierReview model.CVReview
	Reviews []model.CVReview
	CVList map[string]model.CVObject
}

type CVHistoryInfo struct {
	Index int
	CVID string
	CVObject *model.CVObject
}
