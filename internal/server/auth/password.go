package auth

import (
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"golang.org/x/crypto/bcrypt"
)

var _ contract.IPassHasher = &PasswordService{}

// PasswordService - Password service
type PasswordService struct{}

// HashPassword - Hash and salt password
func (*PasswordService) HashPassword(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword - Compare hashed password with plain password
func (*PasswordService) ComparePassword(
	hashedPwd string,
	plainPwd []byte,
) (bool, error) {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false, err
	}

	return true, nil
}
