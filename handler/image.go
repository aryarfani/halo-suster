package handler

import (
	"eniqilo-store/service"

	"github.com/gofiber/fiber/v2"
)

const (
	maxFileSize = 2 * 1024 * 1024 // 2MB
	minFileSize = 10 * 1024       // 10KB
)

func UploadImage(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("File is required")
	}

	// Validate file size
	if fileHeader.Size > maxFileSize || fileHeader.Size < minFileSize {
		return c.Status(fiber.StatusBadRequest).SendString("File size must be between 10KB and 2MB")
	}

	// Validate file extension
	if !service.IsValidImage(fileHeader) {
		return c.Status(fiber.StatusBadRequest).SendString("Only .jpg and .jpeg formats are allowed")
	}

	// Upload File
	err = service.UploadImage(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Record created successfully",
	})
}
