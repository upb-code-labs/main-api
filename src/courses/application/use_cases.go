package application

import "github.com/UPB-Code-Labs/main-api/src/courses/domain/definitions"

type CoursesUseCases struct {
	Repository              definitions.CoursesRepository
	InvitationCodeGenerator definitions.InvitationCodeGenerator
}
