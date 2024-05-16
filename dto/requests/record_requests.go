package requests

type CreateRecordRequest struct {
	IdentityNumber int    `json:"identityNumber" validate:"required,xIntLen=16"`
	Symptoms       string `json:"symptoms" validate:"required,min=1,max=2000"`
	Medications    string `json:"medications" validate:"required,min=1,max=2000"`
}
