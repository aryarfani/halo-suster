package requests

type RegisterUserNurseRequest struct {
	NIP                 string `json:"nip" validate:"required,min=13,max=13,startswith=615"`
	Name                string `json:"name" validate:"required,min=5,max=50"`
	IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,xImageUrl"`
}

type LoginUserNurseRequest struct {
	NIP      string `json:"nip" validate:"required,min=13,max=13,startswith=615"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type UpdatePasswordUserNurseRequest struct {
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UpdateUserNurseRequest struct {
	NIP  string `json:"nip" validate:"required,min=13,max=13,startswith=615"`
	Name string `json:"name" validate:"required,min=5,max=50"`
}
