package model

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID                  uuid.UUID `json:"-"`
	IdentityNumber      int       `json:"identityNumber"`
	PhoneNumber         string    `json:"phoneNumber"`
	Name                string    `json:"name"`
	BirthDate           time.Time `json:"birthDate"`
	Gender              string    `json:"gender"`
	IdentityCardScanImg string    `json:"identityCardScanImg"`
	CreatedAt           time.Time `json:"createdAt"`
}
