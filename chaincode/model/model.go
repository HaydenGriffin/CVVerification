package model

import "time"

// Actor metadata used for an admin and a consumer
type Actor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Available actor type
const (
	ActorAttribute = "actor"
	ActorApplicant = "applicant"
	ActorVerifier = "verifier"
	ActorEmployer = "employer"
	ActorAdmin = "admin"
)

type Admin struct {
	Actor
}

type Applicant struct {
	Actor
}

type UserProfile struct {
	Username	string	`json:"Name"`
	CVHistory []string `json:"CVHistory"`
	Ratings map[string] []CVRating
}

type CVObject struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`
	Speciality	string	`json:"Speciality"`
	CV	string	`json:"CV"`
	CVDate	string	`json:"CVDate"`
}

type CVRating struct {
	Name string `json:"Name"`
	Comment string `json:"Comment"`
	Rating int `json:"Rating"`
}

type CVHistoryInfo struct {
	Index int `json:"index"`
	CVHash string    `json:"cvhash"`
	Timestamp        time.Time `json:"timestamp"`
	CVInReview     int      `json:"cvinreview"`
}

// List of object type stored in the ledger
const (
	ObjectTypeAdmin            = "admin"
	ObjectTypeApplicant        = "applicant"
	ObjectTypeCV         	   = "cv"
	ObjectTypeProfile          = "profile"
	ObjectTypeRating           = "rating"
)

// List of available filter for query resources
const (
	ResourcesFilterAll             = "all"
	ResourcesFilterOnlyAvailable   = "only-available"
	ResourcesFilterOnlyUnavailable = "only-unavailable"
)
