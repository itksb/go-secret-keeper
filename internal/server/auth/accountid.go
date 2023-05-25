package auth

import "github.com/google/uuid"

type userID string

// FieldID field name for user ID
const FieldID userID = "uid"

// GenerateUserID - generates unique user id using uuid
func GenerateUserID() string {
	return uuid.NewString()
}
