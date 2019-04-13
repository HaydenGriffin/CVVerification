package model

// Actor metadata used for an admin and a consumer
type Actor struct {
	ID   string `json:"id"`
	Username string `json:"username"`
}

// Available actor type
const (
	ActorAttribute = "actor"
	ActorApplicant = "applicant"
	ActorVerifier = "verifier"
	ActorEmployer = "employer"
	ActorAdmin = "admin"
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
	CVHistory []string `json:"CVHistory"`
	Reviews map[string] []CVReview
	PublicKey string
}

type VerifierProfile struct {
}

type AdminProfile struct {
}

type EmployerProfile struct {
}

type CVObject struct {
	DocType string `json:"DocType"`
	Name string `json:"Name"`
	Speciality string `json:"Speciality"`
	CV string `json:"CV"`
	CVSections map[string] string `json:"CVSections"`
	CVDate	string `json:"CVDate"`
	Status string `json:"CVStatus"`
}

type CVReview struct {
	VerifierID string `json:"Id"`
	Name string `json:"Name"`
	Comment string `json:"Comment"`
	Rating int `json:"Rating"`
}

// List of object type stored in the ledger
const (
	ObjectTypeApplicant = "applicant"
	ObjectTypeVerifier = "verifier"
	ObjectTypeEmployer = "employer"
	ObjectTypeAdmin = "admin"
	ObjectTypeCV = "cv"
)

const (
	CVInDraft = "draft"
	CVInReview = "in-review"
	CVReviewed = "reviewed"
	CVFinalised = "finalised"
	CVSubmitted = "submitted"
	CVSubmittedRated = "submitted-rated"
)