package implementations

import (
	"github.com/jaevor/go-nanoid"
)

type NanoIdInvitationCodeGenerator struct{}

// Singleton instance
var instance *NanoIdInvitationCodeGenerator

func GetNanoIdInvitationCodeGenerator() *NanoIdInvitationCodeGenerator {
	if instance == nil {
		instance = &NanoIdInvitationCodeGenerator{}
	}

	return instance
}

// Methods
func (generator *NanoIdInvitationCodeGenerator) Generate() (string, error) {
	gen, err := nanoid.Standard(9)
	if err != nil {
		return "", err
	}

	code := gen()
	return code, nil
}
