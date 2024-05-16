package handler

import (
	"eniqilo-store/constant"
	"eniqilo-store/db"
	"eniqilo-store/dto/requests"
	"eniqilo-store/model"
	"eniqilo-store/utils"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserIT(c *fiber.Ctx) error {
	var err error
	var req requests.RegisterUserITRequest
	_ = c.BodyParser(&req)
	if err := utils.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// Validate NIP is not duplicate
	var existingId string
	db.DB.QueryRow("SELECT id FROM users WHERE nip = $1", req.NIP).Scan(&existingId)
	if existingId != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"nip": "NIP is already registered",
			},
		})
	}

	// Hash password
	hashedPassword := utils.HashPassword(req.Password)

	// Insert user into database
	var userId uuid.UUID
	err = db.DB.QueryRow("INSERT INTO users (nip, role, name, password) VALUES ($1, $2, $3, $4) RETURNING id",
		req.NIP, constant.IT, req.Name, hashedPassword).Scan(&userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate access token
	token, _ := utils.GenerateToken(&userId)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"userId":      userId,
			"nip":         req.NIP,
			"name":        req.Name,
			"accessToken": token,
		},
	})
}

func LoginUserIT(c *fiber.Ctx) error {
	var req requests.LoginUserITRequest
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
