package application

import (
	blocksDefinitions "github.com/UPB-Code-Labs/main-api/src/blocks/domain/definitions"
	coursesDefinitions "github.com/UPB-Code-Labs/main-api/src/courses/domain/definitions"
	coursesErrors "github.com/UPB-Code-Labs/main-api/src/courses/domain/errors"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
	laboratoriesErrors "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
	languagesDefinitions "github.com/UPB-Code-Labs/main-api/src/languages/domain/definitions"
	rubricsDefinitions "github.com/UPB-Code-Labs/main-api/src/rubrics/domain/definitions"
	rubricsErrors "github.com/UPB-Code-Labs/main-api/src/rubrics/domain/errors"
	staticFilesDefinitions "github.com/UPB-Code-Labs/main-api/src/static-files/domain/definitions"
	staticFilesDTOs "github.com/UPB-Code-Labs/main-api/src/static-files/domain/dtos"
)

type LaboratoriesUseCases struct {
	CoursesRepository      coursesDefinitions.CoursesRepository
	LaboratoriesRepository definitions.LaboratoriesRepository
	RubricsRepository      rubricsDefinitions.RubricsRepository
	LanguagesRepository    languagesDefinitions.LanguagesRepository
	BlocksRepository       blocksDefinitions.BlockRepository
	StaticFilesRepository  staticFilesDefinitions.StaticFilesRepository
}

func (useCases *LaboratoriesUseCases) CreateLaboratory(dto *dtos.CreateLaboratoryDTO) (laboratory *entities.Laboratory, err error) {
	// Check that the teacher owns the course
	ownsCourse, err := useCases.CoursesRepository.DoesTeacherOwnsCourse(dto.TeacherUUID, dto.CourseUUID)
	if err != nil {
		return nil, err
	}

	if !ownsCourse {
		return nil, laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Create the laboratory
	return useCases.LaboratoriesRepository.SaveLaboratory(dto)
}

func (useCases *LaboratoriesUseCases) GetLaboratory(dto *dtos.GetLaboratoryDTO) (laboratory *entities.Laboratory, err error) {
	// Get the laboratory
	laboratory, err = useCases.LaboratoriesRepository.GetLaboratoryByUUID(dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	// Check that the user is enrolled in the course
	isEnrolled, err := useCases.CoursesRepository.IsUserInCourse(dto.UserUUID, laboratory.CourseUUID)
	if err != nil {
		return nil, err
	}

	if !isEnrolled {
		return nil, coursesErrors.UserNotInCourseError{}
	}

	return laboratory, nil
}

func (useCases *LaboratoriesUseCases) GetLaboratoryInformation(dto *dtos.GetLaboratoryDTO) (laboratoryInformation *dtos.LaboratoryDetailsDTO, err error) {
	// Get the laboratory
	laboratoryInformation, err = useCases.LaboratoriesRepository.GetLaboratoryInformationByUUID(dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	// Check the user is enrolled in the course
	isEnrolled, err := useCases.CoursesRepository.IsUserInCourse(
		dto.UserUUID,
		laboratoryInformation.CourseUUID,
	)
	if err != nil {
		return nil, err
	}

	if !isEnrolled {
		return nil, coursesErrors.UserNotInCourseError{}
	}

	return laboratoryInformation, nil
}

func (useCases *LaboratoriesUseCases) UpdateLaboratory(dto *dtos.UpdateLaboratoryDTO) error {
	// Check that the teacher owns the laboratory / course
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(dto.TeacherUUID, dto.LaboratoryUUID)
	if err != nil {
		return err
	}

	if !teacherOwnsLaboratory {
		return laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Check that the teacher owns the rubric
	if dto.RubricUUID != nil {
		teacherOwnsRubric, err := useCases.RubricsRepository.DoesTeacherOwnRubric(dto.TeacherUUID, *dto.RubricUUID)
		if err != nil {
			return err
		}
		if !teacherOwnsRubric {
			return &rubricsErrors.TeacherDoesNotOwnsRubric{}
		}
	}

	// Update the laboratory
	return useCases.LaboratoriesRepository.UpdateLaboratory(dto)
}

func (useCases *LaboratoriesUseCases) CreateMarkdownBlock(dto *dtos.CreateMarkdownBlockDTO) (blockUUID string, err error) {
	// Check that the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(dto.TeacherUUID, dto.LaboratoryUUID)
	if err != nil {
		return "", err
	}

	if !teacherOwnsLaboratory {
		return "", laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Create the block
	return useCases.LaboratoriesRepository.CreateMarkdownBlock(dto.LaboratoryUUID)
}

func (useCases *LaboratoriesUseCases) CreateTestBlock(reqDTO *dtos.CreateTestBlockDTO) (blockUUID string, err error) {
	// Check that the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(reqDTO.TeacherUUID, reqDTO.LaboratoryUUID)
	if err != nil {
		return "", err
	}

	if !teacherOwnsLaboratory {
		return "", laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Check that the language exists
	_, err = useCases.LanguagesRepository.GetByUUID(reqDTO.LanguageUUID)
	if err != nil {
		return "", err
	}

	// Send the file to the static files microservice
	savedArchiveUUID, err := useCases.StaticFilesRepository.SaveArchive(
		&staticFilesDTOs.SaveStaticFileDTO{
			File:     reqDTO.MultipartFile,
			FileType: "test",
		},
	)
	if err != nil {
		return "", err
	}
	reqDTO.TestArchiveUUID = savedArchiveUUID

	// Save the information in the database
	return useCases.LaboratoriesRepository.CreateTestBlock(reqDTO)
}

func (useCases *LaboratoriesUseCases) GetLaboratoryProgress(dto *dtos.GetLaboratoryProgressDTO) (progress *dtos.LaboratoryProgressDTO, err error) {
	// Check that the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(dto.TeacherUUID, dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	if !teacherOwnsLaboratory {
		return nil, laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Get the total test blocks
	totalTestBlocks, err := useCases.LaboratoriesRepository.GetTotalTestBlocks(dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	// Get the students progress
	studentsProgress, err := useCases.LaboratoriesRepository.GetStudentsProgress(dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	return &dtos.LaboratoryProgressDTO{
		TotalTestBlocks:  totalTestBlocks,
		StudentsProgress: studentsProgress,
	}, nil
}
