package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError map[string]string

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
	validate.RegisterValidation("xImageUrl", validateImageURL)
	validate.RegisterValidation("xBool", validateBoolean)
	validate.RegisterValidation("xNumeric", validateNumeric)
	validate.RegisterValidation("xIntLen", validateNumberOfDigit)

	// register tag name to be validated instead using field name
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateRequest(value interface{}) ValidationError {
	err := validate.Struct(value)
	if err != nil {
		return validatorErrorMessage(err)
	}

	return nil
}

func validatorErrorMessage(validationError error) ValidationError {
	errFields := make(map[string]string)
	for _, err := range validationError.(validator.ValidationErrors) {
		errFields[err.Field()] = msgForTag(err)
	}

	return errFields
}

// func validatorErrorMessageArray(validationError error) []string {
// 	var errFields []string

// 	for _, err := range validationError.(validator.ValidationErrors) {
// 		errFields = append(errFields, err.Field()+" "+msgForTag(err))
// 	}

// 	return errFields
// }

func msgForTag(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "Wajib diisi"
	case "email":
		return "Email tidak valid"
	case "hexcolor":
		return "Kode warna tidak valid"
	case "base64rawurl":
		return "Format harus berupa base64 raw url encoding"
	case "min":
		return fmt.Sprintf("Minimal panjang %s karakter", fieldError.Param())
	case "max":
		return fmt.Sprintf("Maksimal panjang %s karakter", fieldError.Param())
	case "startswith":
		return fmt.Sprintf("Harus dimulai dengan %s", fieldError.Param())
	case "eqfield":
		return "Input tidak sama"
	case "numeric":
		return "Hanya boleh berupa angka"
	case "alpha":
		return "Hanya boleh berupa huruf"
	case "oneof":
		return "Hanya boleh berupa salah satu dari: " + fieldError.Param()
	case "required_if":
		return "Wajib diisi jika " + fieldError.Param()
	case "xIntLen":
		return fmt.Sprintf("Hanya boleh berupa angka %s digit", fieldError.Param())
	default:
		return fieldError.Error()
	}
}

// validateImageURL checks if the provided URL is a valid image URL.
func validateImageURL(fl validator.FieldLevel) bool {
	urlString := fl.Field().String()
	if urlString == "" {
		return true // Allow empty URLs if they are not required
	}

	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}

	// Check if the scheme is HTTP or HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if the URL ends with a valid image file extension
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range validExtensions {
		if strings.HasSuffix(strings.ToLower(parsedURL.Path), ext) {
			return true
		}
	}

	return false
}

// validateBoolean checks if the provided value is a boolean.
func validateBoolean(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()

	switch value := value.(type) {
	case bool:
		return true
	case string:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false
		}
		return boolValue
	default:
		return false
	}
}

func validateNumeric(fl validator.FieldLevel) bool {
	// Get the field value
	value := fl.Field().Interface()

	// Check if the value is numeric
	switch value := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	case string:
		// Attempt to parse string value to float
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	default:
		return false
	}
}

func validateNumberOfDigit(fl validator.FieldLevel) bool {
	field := fl.Field()
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		panic(err.Error())
	}

	v := field.Int()
	if v < 0 {
		panic("negative number")
	}

	n := 0
	for ; v > 0; v /= 10 {
		n += 1
	}

	return n == param
}
