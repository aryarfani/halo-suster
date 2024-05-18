package requests

type RegisterUserITRequest struct {
	NIP      int    `json:"nip" validate:"required,xIntStartsWith=615,nip_genderdigit,nip_validyear,nip_validmonth,nip_validrandomdigits"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoginUserITRequest struct {
	NIP      int    `json:"nip" validate:"required,xIntStartsWith=615,nip_genderdigit,nip_validyear,nip_validmonth,nip_validrandomdigits"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
