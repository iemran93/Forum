package functions

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// hashing password

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// validate user data

func ValidUserData(username, email string) (bool, string) {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,16}$`)
	if emailRegex.MatchString(email) {
		if usernameRegex.MatchString(username) {
			return true, ""
		} else {
			return false, "Username not valid"
		}
	} else {
		return false, "Email not valid"
	}
}
