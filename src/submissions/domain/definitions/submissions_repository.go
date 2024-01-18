package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
)

type SubmissionsRepository interface {
	// SaveSubmission saves the metadata of a new submission in the database
	SaveSubmission(dto *dtos.CreateSubmissionDTO) (submissionUUID string, err error)
	// ResetSubmissionStatus resets the status of a submission to "pending"
	ResetSubmissionStatus(submissionUUID string) (err error)

	// GetStudentSubmission returns the metadata of an student submission
	GetStudentSubmission(studentUUID string, testBlockUUID string) (submission *entities.Submission, err error)
	// GetSubmissionWorkMetadata returns the metadata needed to enqueue a new submission work
	GetSubmissionWorkMetadata(submissionUUID string) (submissionWorkMetadata *entities.SubmissionWork, err error)

	// GetStudentSubmissionArchiveUUIDFromSubmissionUUID returns the UUID of the submission archive in the Static Files Micro-service
	GetStudentSubmissionArchiveUUIDFromSubmissionUUID(submissionUUID string) (archiveUUID string, err error)

	// DoesStudentOwnSubmission returns true if the student owns the submission
	DoesStudentOwnSubmission(studentUUID string, submissionUUID string) (bool, error)

	// GetTeacherOfCourseBySubmissionUUID returns the teacher of the course that the submission belongs to
	GetTeacherOfCourseBySubmissionUUID(submissionUUID string) (teacherUUID string, err error)
}
