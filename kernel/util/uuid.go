package util

import "github.com/google/uuid"

func GetUUID() string {
	uuidObj, _ := uuid.NewV7()
	return uuidObj.String()
}
