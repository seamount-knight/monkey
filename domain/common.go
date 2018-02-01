package domain

import "github.com/satori/go.uuid"

const Source = "001"

func IsValidUUID(validUUID string) bool {
	_, err := uuid.FromString(validUUID)
	return err == nil
}
