package Operations

import (
	"errors"

	"github.com/google/uuid"
)

func validateName(name string) error {
	if name == "" {
		return errors.New("no name inserted")
	} else if len(name) > 60 {
		return errors.New("long name")
	} else if len(name) < 3 {
		return errors.New("short name")
	}
	return nil
}

func validateQuote(quote string) error {
	if quote == "" {
		return errors.New("no quote inserted")
	} else if len(quote) > 300 {
		return errors.New("long quote")
	} else if len(quote) < 10 {
		return errors.New("short quote")
	}
	return nil
}

func ValidateUUID(UUID string) bool {
	_, err := uuid.Parse(UUID)
	return err == nil
}
