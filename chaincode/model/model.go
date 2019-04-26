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
	CVHistory     []string
	Reviews       map[string]map[string][]byte
	PublicReviews map[string][]CVReview
	PublicKey     string
}

type VerifierProfile struct {
	Organisation string
}

type AdminProfile struct {
}

type EmployerProfile struct {
	ProspectiveCVs []string
}

type CVObject struct {
	Name       string
	Date       string
	Industry   string
	Level      string
	CV         string
	CVSections map[string]string
	Status     string
}

type CVReview struct {
	Name         string
	Organisation string
	Type         string
	Comment      string
	Rating       int
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
	CVInDraft   = "draft"
	CVInReview  = "in-review"
	CVSubmitted = "submitted"
	CVWithdrawn = "withdrawn"
)
