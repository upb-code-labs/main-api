package application

import (
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/errors"
	languagesDefinitions "github.com/UPB-Code-Labs/main-api/src/languages/domain/definitions"
)

type BlocksUseCases struct {
	BlocksRepository    definitions.BlockRepository
	LanguagesRepository languagesDefinitions.LanguagesRepository
}

func (useCases *BlocksUseCases) UpdateMarkdownBlockContent(dto dtos.UpdateMarkdownBlockContentDTO) (err error) {
	// Validate the teacher is the owner of the block
	ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsMarkdownBlock(dto.TeacherUUID, dto.BlockUUID)
	if err != nil {
		return err
	}

	if !ownsBlock {
		return errors.TeacherDoesNotOwnBlock{}
	}

	// Update the block
	return useCases.BlocksRepository.UpdateMarkdownBlockContent(dto.BlockUUID, dto.Content)
}

func (useCases *BlocksUseCases) UpdateTestBlock(dto dtos.UpdateTestBlockDTO) (err error) {
	// Validate the teacher is the owner of the block
	ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsTestBlock(dto.TeacherUUID, dto.BlockUUID)
	if err != nil {
		return err
	}

	if !ownsBlock {
		return errors.TeacherDoesNotOwnBlock{}
	}

	// Validate the programming language exists
	_, err = useCases.LanguagesRepository.GetByUUID(dto.LanguageUUID)
	if err != nil {
		return err
	}

	// Overwrite the block's tests archive if the teacher uploaded a new one
	if dto.NewTestArchive != nil {
		// Get the UUID of the block's tests archive
		uuid, err := useCases.BlocksRepository.GetTestArchiveUUIDFromTestBlockUUID(dto.BlockUUID)
		if err != nil {
			return err
		}

		// Send the request to the microservice
		err = useCases.BlocksRepository.OverwriteTestsArchive(uuid, dto.NewTestArchive)
		if err != nil {
			return err
		}
	}

	// Update the block
	return useCases.BlocksRepository.UpdateTestBlock(&dto)
}
