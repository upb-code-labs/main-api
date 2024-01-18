package application

import (
	"mime/multipart"
	"time"

	blocksDefinitions "github.com/UPB-Code-Labs/main-api/src/blocks/domain/definitions"
	laboratoriesDefinitions "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	staticFilesDefinitions "github.com/UPB-Code-Labs/main-api/src/static-files/domain/definitions"
	staticFilesDTOs "github.com/UPB-Code-Labs/main-api/src/static-files/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/errors"
)

type SubmissionUseCases struct {
	StaticFilesRepository   staticFilesDefinitions.StaticFilesRepository
	LaboratoriesRepository  laboratoriesDefinitions.LaboratoriesRepository
	BlocksRepository        blocksDefinitions.BlockRepository
	SubmissionsRepository   definitions.SubmissionsRepository
	SubmissionsQueueManager definitions.SubmissionsQueueManager
}

func (useCases *SubmissionUseCases) CanStudentSubmitToTestBlock(studentUUID string, testBlockUUID string) (bool, error) {
	return useCases.BlocksRepository.CanStudentSubmitToTestBlock(studentUUID, testBlockUUID)
}

func (useCases *SubmissionUseCases) SaveSubmission(dto *dtos.CreateSubmissionDTO) (string, error) {
	// Validate the student can submit to the given test block
	canSubmit, err := useCases.CanStudentSubmitToTestBlock(dto.StudentUUID, dto.TestBlockUUID)
	if err != nil {
		return "", err
	}

	if !canSubmit {
		return "", errors.StudentCannotSubmitToTestBlock{}
	}

	// Validate the laboratory is open
	isLaboratoryOpen, err := useCases.isTestBlockLaboratoryOpen(dto.TestBlockUUID)
	if err != nil {
		return "", err
	}

	if !isLaboratoryOpen {
		return "", errors.LaboratoryIsClosed{}
	}

	// Check if the student already has a submission for the given test block
	previousStudentSubmission, err := useCases.SubmissionsRepository.GetStudentSubmission(dto.StudentUUID, dto.TestBlockUUID)
	if err != nil {
		return "", err
	}

	if previousStudentSubmission != nil {
		// Check if the previous submission was submitted in the last minute
		parsedSubmittedAt, err := time.Parse(time.RFC3339, previousStudentSubmission.SubmittedAt)
		if err != nil {
			return "", err
		}

		if time.Since(parsedSubmittedAt).Minutes() < 1 {
			return "", errors.StudentHasRecentSubmission{}
		}

		// Check if the previous submission is still pending
		finalStatus := "ready"
		if previousStudentSubmission.Status != finalStatus {
			return "", errors.StudentHasPendingSubmission{}
		}

		// If the student already has a submission, reset its status and overwrite the archive
		err = useCases.resetSubmissionStatus(previousStudentSubmission, dto.SubmissionArchive)
		if err != nil {
			return "", err
		}

		// Submit the work to the submissions queue
		err = useCases.submitWorkToQueue(previousStudentSubmission.UUID)
		if err != nil {
			return "", err
		}

		return previousStudentSubmission.UUID, nil
	} else {
		// If the student doesn't have a submission, create a new one
		submissionUUID, err := useCases.createSubmission(dto)
		if err != nil {
			return "", err
		}

		// Submit the work to the submissions queue
		err = useCases.submitWorkToQueue(submissionUUID)
		if err != nil {
			return "", err
		}

		return submissionUUID, nil
	}
}

func (useCases *SubmissionUseCases) isTestBlockLaboratoryOpen(testBlockUUID string) (bool, error) {
	// Get the UUID of the laboratory the test block belongs to
	laboratoryUUID, err := useCases.BlocksRepository.GetTestBlockLaboratoryUUID(testBlockUUID)
	if err != nil {
		return false, err
	}

	// Get the laboratory
	laboratory, err := useCases.LaboratoriesRepository.GetLaboratoryInformationByUUID(laboratoryUUID)
	if err != nil {
		return false, err
	}

	// Check if the laboratory is open
	parsedClosingDate, err := time.Parse(time.RFC3339, laboratory.DueDate)
	if err != nil {
		return false, err
	}

	if time.Now().After(parsedClosingDate) {
		return false, nil
	}

	return true, nil
}

func (useCases *SubmissionUseCases) resetSubmissionStatus(previousStudentSubmission *entities.Submission, newArchive *multipart.File) error {
	// Get the UUID of the .zip archive in the static files microservice
	archiveUUID, err := useCases.SubmissionsRepository.GetStudentSubmissionArchiveUUIDFromSubmissionUUID(previousStudentSubmission.UUID)
	if err != nil {
		return err
	}

	// Overwrite the archive in the static files microservice
	err = useCases.StaticFilesRepository.OverwriteArchive(
		&staticFilesDTOs.OverwriteStaticFileDTO{
			FileUUID: archiveUUID,
			FileType: "submission",
			File:     newArchive,
		},
	)
	if err != nil {
		return err
	}

	// Reset the submission status
	err = useCases.SubmissionsRepository.ResetSubmissionStatus(previousStudentSubmission.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *SubmissionUseCases) createSubmission(dto *dtos.CreateSubmissionDTO) (string, error) {
	// Save the .zip archive in the static files microservice
	archiveUUID, err := useCases.StaticFilesRepository.SaveArchive(
		&staticFilesDTOs.SaveStaticFileDTO{
			FileType: "submission",
			File:     dto.SubmissionArchive,
		},
	)
	if err != nil {
		return "", err
	}

	dto.SavedArchiveUUID = archiveUUID

	// Save the submission
	submissionUUID, err := useCases.SubmissionsRepository.SaveSubmission(dto)
	if err != nil {
		return "", err
	}

	return submissionUUID, nil
}

func (useCases *SubmissionUseCases) submitWorkToQueue(submissionUUID string) error {
	// Get the submission work
	submissionWork, err := useCases.SubmissionsRepository.GetSubmissionWorkMetadata(submissionUUID)
	if err != nil {
		return err
	}

	// Send the submission work to the submissions queue
	err = useCases.SubmissionsQueueManager.QueueWork(submissionWork)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *SubmissionUseCases) GetSubmissionStatus(studentUUID, testBlockUUID string) (*dtos.SubmissionStatusUpdateDTO, error) {
	// Check if the student could submit to the given test block
	canSubmit, err := useCases.CanStudentSubmitToTestBlock(studentUUID, testBlockUUID)
	if err != nil {
		return nil, err
	}

	if !canSubmit {
		return nil, errors.StudentCannotSubmitToTestBlock{}
	}

	// Get the submission
	submission, err := useCases.SubmissionsRepository.GetStudentSubmission(studentUUID, testBlockUUID)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, errors.StudentSubmissionNotFound{}
	}

	// Get the submission status
	dto := dtos.SubmissionStatusUpdateDTO{
		SubmissionUUID:   submission.UUID,
		SubmissionStatus: submission.Status,
		TestsPassed:      submission.Passing,
		TestsOutput:      submission.Stdout,
	}

	return &dto, nil
}

// GetSubmissionArchive Use case to return the bytes of the `zip` archive of a submission
func (useCases *SubmissionUseCases) GetSubmissionArchive(dto *dtos.GetSubmissionArchiveDTO) ([]byte, error) {
	// Check if the user has access to the submission
	if dto.UserRole == "teacher" {
		// If the user is a teacher, check if is the teacher of the course that the submission belongs to
		teacherUUID, err := useCases.SubmissionsRepository.GetTeacherOfCourseBySubmissionUUID(dto.SubmissionUUID)
		if err != nil {
			return nil, err
		}

		if teacherUUID != dto.UserUUID {
			return nil, errors.UserDoesNotHaveAccessToSubmission{}
		}
	} else {
		// If the user is a student, check if the student owns the submission
		ownsSubmission, err := useCases.SubmissionsRepository.DoesStudentOwnSubmission(dto.UserUUID, dto.SubmissionUUID)
		if err != nil {
			return nil, err
		}

		if !ownsSubmission {
			return nil, errors.UserDoesNotHaveAccessToSubmission{}
		}
	}

	// Get the UUID of the .zip archive in the static files microservice
	archiveUUID, err := useCases.SubmissionsRepository.GetStudentSubmissionArchiveUUIDFromSubmissionUUID(dto.SubmissionUUID)
	if err != nil {
		return nil, err
	}

	// Get the bytes of the .zip archive
	archiveBytes, err := useCases.StaticFilesRepository.GetArchiveBytes(&staticFilesDTOs.StaticFileArchiveDTO{
		FileUUID: archiveUUID,
		FileType: "submission",
	})
	if err != nil {
		return nil, err
	}

	return archiveBytes, nil
}
