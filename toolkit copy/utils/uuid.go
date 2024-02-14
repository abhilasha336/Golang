package utils

import "github.com/google/uuid"

func IsValidUUID(input string) bool {
	_, err := uuid.Parse(input)
	return err == nil
}
