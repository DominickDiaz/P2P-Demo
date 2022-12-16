package sec

import (
	"github.com/google/uuid"
)

// GenerateToken  generates a new token using the UUID package.
func GenerateToken() string {
	// Generate a new UUID.
	uid, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	// Return the UUID as a string.
	return uid.String()
}
