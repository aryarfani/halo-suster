package handler

import (
	"eniqilo-store/constant"
	"eniqilo-store/db"
	"eniqilo-store/dto/requests"
	"eniqilo-store/model"
	"eniqilo-store/utils"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 5)
	offset := c.QueryInt("offset", 0)
	name := c.Query("name")
	role := c.Query("role")
	nip := c.QueryInt("nip")
	createdAt := c.Query("created_at")
	userId := c.Query("userId")

	statement := "SELECT id, nip, name, created_at FROM users WHERE deleted_at IS NULL"

	if userId != "" {
		statement += " AND id = '" + userId + "' "
	}
	if name != "" {
		statement += " AND name ILIKE '%" + name + "%'"
	}

	if nip != 0 {
		statement += fmt.Sprintf(" AND nip::TEXT LIKE '%%%d%%'", nip)
	}

	if createdAt == "asc" || createdAt == "desc" {
		statement += fmt.Sprintf(" ORDER BY created_at %s ", createdAt)
	}

	if role == "it" || role == "nurse" {
		statement += fmt.Sprintf(" AND role = '%s' ", role)
	}

	statement += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	log.Println(statement)
	rows, err := db.DB.Query(statement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		rows.Scan(&user.ID, &user.NIP, &user.Name, &user.CreatedAt)
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    users,
	})
}

func RegisterUserNurse(c *fiber.Ctx) error {
	var err error
	var req requests.RegisterUserNurseRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// Validate NIP is not duplicate
	var existingId string
	db.DB.QueryRow("SELECT id FROM users WHERE nip = $1 AND deleted_at IS NULL", req.NIP).Scan(&existingId)
	if existingId != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"nip": "NIP is already registered",
			},
		})
	}

	// Insert user into database
	var userId uuid.UUID
	err = db.DB.QueryRow("INSERT INTO users (nip, role, name, identity_card_scan_img) VALUES ($1, $2, $3, $4) RETURNING id",
		req.NIP, constant.Nurse, req.Name, req.IdentityCardScanImg).Scan(&userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"userId": userId,
			"nip":    req.NIP,
			"name":   req.Name,
		},
	})
}

func LoginUserNurse(c *fiber.Ctx) error {
	var req requests.LoginUserNurseRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// Get user by nip
	var user model.User
	err := db.DB.QueryRow("SELECT id, nip, name, password FROM users WHERE nip = $1", req.NIP).
		Scan(&user.ID, &user.NIP, &user.Name, &user.Password)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Validate password is correct
	if !utils.ComparePassword(user.Password, req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	// Generate access token
	token, _ := utils.GenerateToken(&user.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"userId":      user.ID,
			"nip":         user.NIP,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}

func ChangePasswordUserNurse(c *fiber.Ctx) error {
	var req requests.UpdatePasswordUserNurseRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	userId := c.Params("id")

	// Get user by id
	var user model.User
	err := db.DB.QueryRow("SELECT id, nip, role FROM users WHERE id = $1", userId).
		Scan(&user.ID, &user.NIP, &user.Role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update password
	hashedPassword := utils.HashPassword(req.Password)
	_, err = db.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", hashedPassword, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

func UpdateUserNurse(c *fiber.Ctx) error {
	userId := c.Params("id")

	//* Validate user exists
	var doesUserExist bool
	err := db.DB.QueryRow("SELECT EXISTS (SELECT * FROM users WHERE id = $1 AND role = $2 AND deleted_at IS NULL)", userId, constant.Nurse).
		Scan(&doesUserExist)
	if !doesUserExist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var req requests.UpdateUserNurseRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	//* Validate NIP is unique
	var doesNipExist bool
	_ = db.DB.QueryRow("SELECT EXISTS (SELECT * FROM users WHERE nip = $1 AND id <> $2)", req.NIP, userId).
		Scan(&doesNipExist)
	if doesNipExist {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "NIP is already registered",
		})
	}

	// Update nurse
	_, err = db.DB.Exec("UPDATE users SET nip = $1, name = $2 WHERE id = $3", req.NIP, req.Name, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Nurse updated successfully",
	})
}

func DeleteUserNurse(c *fiber.Ctx) error {
	userId := c.Params("id")

	//* Validate user exists
	var doesUserExist bool
	err := db.DB.QueryRow("SELECT EXISTS (SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL AND role = $2)", userId, constant.Nurse).
		Scan(&doesUserExist)
	if !doesUserExist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Soft Delete User
	_, err = db.DB.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1", userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Nurse deleted successfully",
	})
}
