package requests

type RegisterUserNurseRequest struct {
	NIP                 int    `json:"nip" validate:"required,xIntStartsWith=303,nip_genderdigit,nip_validyear,nip_validmonth,nip_validrandomdigits"`
	Name                string `json:"name" validate:"required,min=5,max=50"`
	IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,xImageUrl"`
}

type LoginUserNurseRequest struct {
	NIP      int    `json:"nip" validate:"required,xIntStartsWith=303,nip_genderdigit,nip_validyear,nip_validmonth,nip_validrandomdigits"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type UpdatePasswordUserNurseRequest struct {
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UpdateUserNurseRequest struct {
	NIP  int    `json:"nip" validate:"required,xIntStartsWith=303,nip_genderdigit,nip_validyear,nip_validmonth,nip_validrandomdigits"`
	Name string `json:"name" validate:"required,min=5,max=50"`
}
