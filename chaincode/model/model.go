package model

// Actor metadata used for an admin and a consumer
type Actor struct {
	ID       string
	Username string
}

// Available actor type
const (
	ActorAttribute = "actor"
	ActorApplicant = "applicant"
	ActorVerifier  = "verifier"
	ActorEmployer  = "employer"
	ActorAdmin     = "admin"
)

type Applicant struct {
	Actor
	Profile ApplicantProfile
}

type Verifier struct {
	Actor
	Profile VerifierProfile
}

type Employer struct {
	Actor
	Profile EmployerProfile
}

type Admin struct {
	Actor
	Profile AdminProfile
}

type ApplicantProfile struct {
	CVHistory []string
	Reviews   map[string]map[string][]byte
	PublicKey string
}

type VerifierProfile struct {
}

type AdminProfile struct {
}

type EmployerProfile struct {
}

type CVObject struct {
	Name       string
	Date     string
	Industry   string
	Level      string
	CV         string
	CVSections map[string]string
	Status     string
}

type CVReview struct {
	VerifierID string
	Name       string
	Comment    string
	Rating     int
}

// List of object type stored in the ledger
const (
	ObjectTypeApplicant = "applicant"
	ObjectTypeVerifier  = "verifier"
	ObjectTypeEmployer  = "employer"
	ObjectTypeAdmin     = "admin"
	ObjectTypeCV        = "cv"
)

const (
	CVInDraft        = "draft"
	CVInReview       = "in-review"
	CVReviewed       = "reviewed"
	CVFinalised      = "finalised"
	CVSubmitted      = "submitted"
	CVSubmittedRated = "submitted-rated"
)
