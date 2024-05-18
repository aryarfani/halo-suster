package handler

import (
	"eniqilo-store/config/storage"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	if !IsValidImage(fileHeader) {
		return c.Status(fiber.StatusBadRequest).SendString("Only .jpg and .jpeg formats are allowed")
	}

	// Upload File
	newFileName := uuid.NewString() + ".jpg"
	uploadOutput, err := storage.Upload(newFileName, fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded sucessfully",
		"data": fiber.Map{
			"imageUrl": uploadOutput.Location,
		},
	})
}

func IsValidImage(fileHeader *multipart.FileHeader) bool {
	ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
	return ext == "jpg" || ext == "jpeg"
}
