package model

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	ID                    uuid.UUID `json:"id"`
	PatientIdentityNumber string    `json:"patient_identity_number"`
	Symptoms              string    `json:"symptoms"`
	Medications           string    `json:"medications"`
	UserId                string    `json:"user_id"`
	CreatedAt             time.Time `json:"createdAt"`
}
