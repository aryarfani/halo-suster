package handler

import (
	"eniqilo-store/db"
	"eniqilo-store/dto/requests"
	"eniqilo-store/model"
	"eniqilo-store/utils"
	"fmt"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

func GetPatients(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 5)
	offset := c.QueryInt("offset", 0)
	name := c.Query("name")
	identityNumber := c.QueryInt("identity_number")
	phoneNumber := c.Query("phone_number")
	createdAt := c.Query("created_at")

	statement := "SELECT identity_number, phone_number, name, birth_date, gender, identity_card_scan_img, created_at FROM patients "

	if name != "" {
		statement += " AND name ILIKE '%" + name + "%'"
	}

	if identityNumber != 0 {
		statement += fmt.Sprintf(" AND nip = %d", identityNumber)
	}

	if phoneNumber != "" {
		statement += " AND phone_number ILIKE '%" + phoneNumber + "%'"
	}

	if createdAt == "asc" || createdAt == "desc" {
		statement += fmt.Sprintf(" ORDER BY created_at %s ", createdAt)
	}

	statement += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.DB.Query(statement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	patients := []model.Patient{}
	for rows.Next() {
		var patient model.Patient
		if err := rows.Scan(&patient.IdentityNumber, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &patient.IdentityCardScanImg, &patient.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		patients = append(patients, patient)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    patients,
	})
}

func CreatePatient(c *fiber.Ctx) error {
	var err error
	var req requests.CreatePatientRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// Validate NIP is not duplicate
	var existingId string
	db.DB.QueryRow("SELECT id FROM patients WHERE identity_number = $1", req.IdentityNumber).Scan(&existingId)
	if existingId != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"nip": "NIP is already registered",
			},
		})
	}

	// Insert patient into database
	var patientId uuid.UUID
	err = db.DB.QueryRow("INSERT INTO patients (identity_number, phone_number, name, birth_date, gender, identity_card_scan_img) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.IdentityNumber, req.PhoneNumber, req.Name, req.BirthDate, req.Gender, req.IdentityCardScanImg).Scan(&patientId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Patient registered successfully",
	})
}
