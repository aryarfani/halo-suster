package service

import (
	"eniqilo-store/config/storage"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IsValidImage(fileHeader *multipart.FileHeader) bool {
	ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
	return ext == "jpg" || ext == "jpeg"
}

func UploadImage(fileHeader *multipart.FileHeader) error {
	f, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the file into a byte slice
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	newFileName := uuid.NewString() + ".jpg"
	err = storage.Storage.Set(newFileName, fileBytes, time.Duration(time.Now().Nanosecond()))
	if err != nil {
		return err
	}

	return nil
}
