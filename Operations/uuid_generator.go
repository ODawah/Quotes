package Operations

import (
	"github.com/google/uuid"
)

func UuidGenerator() string {
	id := uuid.New().String()
	return id
}
