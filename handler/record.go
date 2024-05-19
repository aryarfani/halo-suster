package handler

import (
	"eniqilo-store/db"
	"eniqilo-store/dto/requests"
	"eniqilo-store/dto/responses"
	"eniqilo-store/utils"
	"fmt"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

func GetRecords(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 5)
	offset := c.QueryInt("offset", 0)

	identityNumber := c.QueryInt("identityDetail.identityNumber")
	userId := c.Query("createdBy.userId")
	nip := c.QueryInt("createdBy.nip")
	createdAt := c.Query("created_at")

	statement := `
	SELECT
		p.identity_number as patient_identity_number,
		p.phone_number as patient_phone_number,
		p.name as patient_name,
		p.birth_date as patient_birth_date,
		p.gender as patient_gender,
		p.identity_card_scan_img as patient_identity_card_scan_img,
		r.symptoms,
		r.medications,
		r.created_at,
		u.nip as user_nip,
		u.name as user_name,
		u.id as user_id
	FROM
		records as r
		JOIN patients as p ON r.patient_identity_number = p.identity_number
		JOIN users as u ON r.user_id = u.id
	WHERE
		1 = 1
	`

	if identityNumber != 0 {
		statement += fmt.Sprintf(" AND p.identity_number = %d", identityNumber)
	}
	_, err := uuid.Parse(userId)
	if userId != "" && err == nil {
		statement += fmt.Sprintf(" AND u.id = '%s'", userId)
	}

	if nip != 0 {
		statement += fmt.Sprintf(" AND u.nip = %d", nip)
	}

	if createdAt == "asc" || createdAt == "desc" {
		statement += fmt.Sprintf(" ORDER BY created_at %s ", createdAt)
	}

	statement += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	fmt.Println(statement)
	rows, err := db.DB.Query(statement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	records := []responses.PatientRecordResponse{}
	for rows.Next() {
		var record responses.PatientRecordResponse
		if err := rows.Scan(
			&record.IdentityDetail.IdentityNumber,
			&record.IdentityDetail.PhoneNumber,
			&record.IdentityDetail.Name,
			&record.IdentityDetail.BirthDate,
			&record.IdentityDetail.Gender,
			&record.IdentityDetail.IdentityCardScanImg,
			&record.Symptoms,
			&record.Medications,
			&record.CreatedAt,
			&record.CreatedBy.NIP,
			&record.CreatedBy.Name,
			&record.CreatedBy.UserID,
		); err != nil {
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
