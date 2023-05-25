package auth

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordHashAndSalt(t *testing.T) {
	password := []byte("password123")
	ps := PasswordService{}
	hashedPassword, _ := ps.HashPassword(password)

	// Проверяем, что хешированный пароль не пустой
	if hashedPassword == "" {
		t.Errorf("Hashed password is empty")
	}

	// Проверяем, что хешированный пароль можно верифицировать
	match, err := ps.ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Failed to compare passwords: %v", err)
	}
	if !match {
		t.Errorf("Password verification failed")
	}
}

func TestPasswordCompare(t *testing.T) {
	plainPassword := []byte("password123")
	ps := PasswordService{}
	hashedPassword, _ := bcrypt.GenerateFromPassword(plainPassword, bcrypt.DefaultCost)

	// Проверяем, что верификация с верным паролем проходит успешно
	match, err := ps.ComparePassword(string(hashedPassword), plainPassword)
	if err != nil {
		t.Errorf("Failed to compare passwords: %v", err)
	}
	if !match {
		t.Errorf("Password verification failed")
	}

	// Проверяем, что верификация с неверным паролем возвращает ошибку
	invalidPassword := []byte("wrongpassword")
	match, err = ps.ComparePassword(string(hashedPassword), invalidPassword)
	if err == nil {
		t.Errorf("Expected password verification to fail, but it succeeded")
	}
	if match {
		t.Errorf("Unexpected password verification success")
	}
}
