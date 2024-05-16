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

func GetRecords(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 5)
	offset := c.QueryInt("offset", 0)

	identityNumber := c.QueryInt("identity_number")
	userId := c.Query("user_id")
	// nip := c.QueryInt("nip")
	createdAt := c.Query("created_at")

	statement := "SELECT id, patient_identity_number, symptoms, medications, user_id, created_at FROM records WHERE 1 = 1"

	if identityNumber != 0 {
		statement += fmt.Sprintf(" AND nip = %d", identityNumber)
	}

	if userId != "" {
		statement += fmt.Sprintf(" AND user_id = %s", userId)
	}

	// if nip != "" {
	// 	statement += fmt.Sprintf(" AND user_id = %d", nip)
	// }

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

	records := []model.Record{}
	for rows.Next() {
		var record model.Record
		if err := rows.Scan(&record.ID, &record.PatientIdentityNumber, &record.Symptoms, &record.Medications, &record.UserId, &record.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		records = append(records, record)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    records,
	})
}

func CreateRecord(c *fiber.Ctx) error {
	var err error
	var req requests.CreateRecordRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// Validate [patient] with identityNumber exists
	var existingId string
	db.DB.QueryRow("SELECT id FROM patients WHERE identity_number = $1", req.IdentityNumber).Scan(&existingId)
	if existingId == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	// Insert record into database
	userId := c.Locals("user_id")
	var recordId uuid.UUID
	err = db.DB.QueryRow("INSERT INTO records (patient_identity_number, symptoms, medications, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		req.IdentityNumber, req.Symptoms, req.Medications, userId).Scan(&recordId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Record created successfully",
	})
}
