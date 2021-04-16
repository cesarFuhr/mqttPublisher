package ports

import (
	"time"

	"github.com/cesarFuhr/validator"
)

var (
	scopeV         = validator.NewStringValidator("scope", true, validator.StrLength(1, 50))
	expirationV    = validator.NewStringValidator("expiration", true, validator.StrDate(time.RFC3339))
	keyIDV         = validator.NewStringValidator("keyID", true, validator.StrUUID())
	dataV          = validator.NewStringValidator("data", true, validator.StrLength(1, 1000))
	encryptedDataV = validator.NewStringValidator("encryptedData", true, validator.StrLength(1, 4000))
)

type keysValidator struct{}

func (v keysValidator) GetValidator(keyID string) error {
	if err := keyIDV.Validate(keyID); err != nil {
		return err
	}
	return nil
}

func (v keysValidator) FindValidator(scope string) error {
	if err := scopeV.Validate(scope); err != nil {
		return err
	}
	return nil
}

type encryptValidator struct{}
