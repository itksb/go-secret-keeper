package auth

import (
	"github.com/google/uuid"
	"testing"
)

func TestGenerateUserID(t *testing.T) {
	userID := GenerateUserID()

	// Check if the generated user ID is a valid UUID
	_, err := uuid.Parse(userID)
	if err != nil {
		t.Errorf("Generated user ID is not a valid UUID: %s", err)
	}
}

func TestGenerateUniqueUserIDs(t *testing.T) {
	// Generate multiple user IDs and ensure they are all unique
	userIDs := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		generatedUserID := GenerateUserID()

		if userIDs[generatedUserID] {
			t.Errorf("Duplicate user ID generated: %s", generatedUserID)
		}

		userIDs[generatedUserID] = true
	}
}
