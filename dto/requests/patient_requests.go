package requests

type CreatePatientRequest struct {
	IdentityNumber      int    `json:"identityNumber" validate:"required,xIntLen=16"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,startswith=+62,min=10,max=15"`
	Name                string `json:"name" validate:"required,min=3,max=30"`
	BirthDate           string `json:"birthDate" validate:"required"`
	Gender              string `json:"gender" validate:"required,oneof=male female"`
	IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,xImageUrl"`
}
