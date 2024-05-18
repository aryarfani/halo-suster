package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	validate.RegisterValidation("xIntStartsWith", validateIntStartsWith)
	validate.RegisterValidation("nip_genderdigit", validGenderDigit)
	validate.RegisterValidation("nip_validyear", validYear)
	validate.RegisterValidation("nip_validmonth", validMonth)
	validate.RegisterValidation("nip_validrandomdigits", validRandomDigits)

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

func validateIntStartsWith(fl validator.FieldLevel) bool {
	// Get the field's value as an integer
	value := fl.Field().Int()

	// Get the parameter (the starting digit) from the tag
	param := fl.Param()

	// Convert the integer value to a string
	valueStr := strconv.FormatInt(value, 10)

	// Check if the value starts with the specified digit
	return strings.HasPrefix(valueStr, param)
}

// Custom validation function to check the fourth digit for gender
func validGenderDigit(fl validator.FieldLevel) bool {
	nip := strconv.FormatInt(fl.Field().Int(), 10)
	if len(nip) < 4 {
		return false
	}
	genderDigit := nip[3]
	return genderDigit == '1' || genderDigit == '2'
}

// Custom validation function to check if year is between 2000 and current year
func validYear(fl validator.FieldLevel) bool {
	nip := strconv.FormatInt(fl.Field().Int(), 10)
	if len(nip) < 8 {
		return false
	}
	year, err := strconv.Atoi(nip[4:8])
	if err != nil {
		return false
	}
	currentYear := time.Now().Year()
	return year >= 2000 && year <= currentYear
}

// Custom validation function to check if month is between 01 and 12
func validMonth(fl validator.FieldLevel) bool {
	nip := strconv.FormatInt(fl.Field().Int(), 10)
	if len(nip) < 10 {
		return false
	}
	month, err := strconv.Atoi(nip[8:10])
	if err != nil {
		return false
	}
	return month >= 1 && month <= 12
}

// Custom validation function to check if the last digits are in valid range
func validRandomDigits(fl validator.FieldLevel) bool {
	nip := strconv.FormatInt(fl.Field().Int(), 10)
	if len(nip) < 11 {
		return false
	}
	randomDigits := nip[10:]
	return len(randomDigits) >= 3 && len(randomDigits) <= 5
}
