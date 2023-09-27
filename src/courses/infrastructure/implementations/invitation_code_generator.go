package implementations

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
func (generator *NanoIdInvitationCodeGenerator) Generate() string {
	return ""
}
