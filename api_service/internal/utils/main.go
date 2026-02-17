package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"reflect"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hash), err
}
func CheckPassword(rawPassword string, encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(rawPassword)) == nil
}

var Validator = validator.New()

func ValidationErrorsToMap(err error, obj interface{}) map[string]string {
	errorsMap := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return errorsMap
	}

	t := reflect.TypeOf(obj).Elem()

	for _, e := range validationErrors {
		field, _ := t.FieldByName(e.StructField())
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = e.Field()
		}

		msg := field.Tag.Get(e.Tag() + "_msg")
		if msg != "" {
			errorsMap[fieldName] = msg
			continue
		}

		errorsMap[fieldName] = defaultValidationMessage(e)
	}

	return errorsMap
}
func defaultValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Minimum length is " + e.Param()
	case "max":
		return "Maximum length is " + e.Param()
	default:
		return "Invalid value"
	}
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range result {
		n, _ := rand.Int(rand.Reader, charsetLen)
		result[i] = charset[n.Int64()]
	}

	return string(result)
}
