package utils

import (
	"eniqilo-store/config"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	num, err := strconv.Atoi(config.GetConfig("BCRYPT_SALT"))
	if err != nil {
		num = bcrypt.DefaultCost
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), num)

	return string(hashed)
}

func ComparePassword(hashedPassword string, password string) bool {
	ok := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return ok == nil
}
