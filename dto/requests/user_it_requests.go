package requests

type RegisterUserITRequest struct {
	NIP      string `json:"nip" validate:"required,min=13,max=13,startswith=615"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoginUserITRequest struct {
	NIP      string `json:"nip" validate:"required,min=13,max=13,startswith=615"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
