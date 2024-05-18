package responses

type IdentityDetail struct {
	IdentityNumber      string `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

type CreatedBy struct {
	NIP    string `json:"nip"`
	Name   string `json:"name"`
	UserID string `json:"userId"`
}

type PatientRecordResponse struct {
	IdentityDetail IdentityDetail `json:"identityDetail"`
	Symptoms       string         `json:"symptoms"`
	Medications    string         `json:"medications"`
	CreatedAt      string         `json:"createdAt"`
	CreatedBy      CreatedBy      `json:"createdBy"`
}
