package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID    `json:"userId"`
	NIP                 uint         `json:"nip"`
	Name                string       `json:"name"`
	Role                string       `json:"-"`
	Password            string       `json:"-"`
	IdentityCardScanImg string       `json:"identityCardScanImg,omitempty"`
	DeletedAt           sql.NullTime `json:"-"`
	CreatedAt           time.Time    `json:"createdAt"`
}
