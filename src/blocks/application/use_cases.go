package application

import (
	"errors"

	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	blocksErrors "github.com/UPB-Code-Labs/main-api/src/blocks/domain/errors"
	languagesDefinitions "github.com/UPB-Code-Labs/main-api/src/languages/domain/definitions"
	staticFilesDefinitions "github.com/UPB-Code-Labs/main-api/src/static-files/domain/definitions"
	staticFilesDTOs "github.com/UPB-Code-Labs/main-api/src/static-files/domain/dtos"
)

type BlocksUseCases struct {
	BlocksRepository      definitions.BlockRepository
	LanguagesRepository   languagesDefinitions.LanguagesRepository
	StaticFilesRepository staticFilesDefinitions.StaticFilesRepository
}

func (useCases *BlocksUseCases) UpdateMarkdownBlockContent(dto dtos.UpdateMarkdownBlockContentDTO) (err error) {
	// Validate the teacher is the owner of the block
	ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsMarkdownBlock(dto.TeacherUUID, dto.BlockUUID)
	if err != nil {
		return err
	}

	if !ownsBlock {
		return blocksErrors.TeacherDoesNotOwnBlock{}
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
		return blocksErrors.TeacherDoesNotOwnBlock{}
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
		err = useCases.StaticFilesRepository.OverwriteArchive(
			&staticFilesDTOs.OverwriteStaticFileDTO{
				FileUUID: uuid,
				FileType: "test",
				File:     dto.NewTestArchive,
			},
		)
		if err != nil {
			return err
		}
	}

	// Update the block
	return useCases.BlocksRepository.UpdateTestBlock(&dto)
}

func (useCases *BlocksUseCases) DeleteMarkdownBlock(dto dtos.DeleteBlockDTO) (err error) {
	// Validate the teacher is the owner of the block
	ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsMarkdownBlock(dto.TeacherUUID, dto.BlockUUID)
	if err != nil {
		return err
	}

	if !ownsBlock {
		return blocksErrors.TeacherDoesNotOwnBlock{}
	}

	// Delete the block
	return useCases.BlocksRepository.DeleteMarkdownBlock(dto.BlockUUID)
}

func (useCases *BlocksUseCases) DeleteTestBlock(dto dtos.DeleteBlockDTO) (err error) {
	// Validate the teacher is the owner of the block
	ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsTestBlock(dto.TeacherUUID, dto.BlockUUID)
	if err != nil {
		return err
	}

	if !ownsBlock {
		return blocksErrors.TeacherDoesNotOwnBlock{}
	}

	// Delete the block
	return useCases.BlocksRepository.DeleteTestBlock(dto.BlockUUID)
}

func (useCases *BlocksUseCases) SwapBlocks(dto dtos.SwapBlocksDTO) (err error) {
	var blockNotFoundError *blocksErrors.BlockNotFound

	// Check if the blocks exists
	firstBlockAsMarkdown, firstBlockAsMarkdownErr := useCases.BlocksRepository.GetMarkdownBlockByUUID(dto.FirstBlockUUID)
	if firstBlockAsMarkdownErr != nil {
		// Forward the error if it's not a `BlockNotFound` error
		if !errors.As(firstBlockAsMarkdownErr, &blockNotFoundError) {
			return firstBlockAsMarkdownErr
		}
	}

	fistBlockAsTest, firstBlockAsTestErr := useCases.BlocksRepository.GetTestBlockByUUID(dto.FirstBlockUUID)
	if firstBlockAsTestErr != nil {
		if !errors.As(firstBlockAsTestErr, &blockNotFoundError) {
			return firstBlockAsTestErr
		}
	}

	// Return an error if the block was not found
	noFirstBlockFound := firstBlockAsMarkdown == nil && fistBlockAsTest == nil
	if noFirstBlockFound {
		return blocksErrors.BlockNotFound{}
	}

	secondBlockAsMarkdown, secondBlockAsMarkdownErr := useCases.BlocksRepository.GetMarkdownBlockByUUID(dto.SecondBlockUUID)
	if secondBlockAsMarkdownErr != nil {
		if !errors.As(secondBlockAsMarkdownErr, &blockNotFoundError) {
			return secondBlockAsMarkdownErr
		}
	}

	secondBlockAsTest, secondBlockAsTestErr := useCases.BlocksRepository.GetTestBlockByUUID(dto.SecondBlockUUID)
	if secondBlockAsTestErr != nil {
		if !errors.As(secondBlockAsTestErr, &blockNotFoundError) {
			return secondBlockAsTestErr
		}
	}

	// Return an error if the block was not found
	noSecondBlockFound := secondBlockAsMarkdown == nil && secondBlockAsTest == nil
	if noSecondBlockFound {
		return blocksErrors.BlockNotFound{}
	}

	// Validate the teacher is the owner of the blocks
	if firstBlockAsMarkdown != nil {
		ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsMarkdownBlock(dto.TeacherUUID, firstBlockAsMarkdown.UUID)
		if err != nil {
			return err
		}

		if !ownsBlock {
			return blocksErrors.TeacherDoesNotOwnBlock{}
		}
	} else {
		ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsTestBlock(dto.TeacherUUID, fistBlockAsTest.UUID)
		if err != nil {
			return err
		}

		if !ownsBlock {
			return blocksErrors.TeacherDoesNotOwnBlock{}
		}
	}

	if secondBlockAsMarkdown != nil {
		ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsMarkdownBlock(dto.TeacherUUID, secondBlockAsMarkdown.UUID)
		if err != nil {
			return err
		}

		if !ownsBlock {
			return blocksErrors.TeacherDoesNotOwnBlock{}
		}
	} else {
		ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsTestBlock(dto.TeacherUUID, secondBlockAsTest.UUID)
		if err != nil {
			return err
		}

		if !ownsBlock {
			return blocksErrors.TeacherDoesNotOwnBlock{}
		}
	}

	// Swap the blocks
	return useCases.BlocksRepository.SwapBlocks(dto.FirstBlockUUID, dto.SecondBlockUUID)
}

// GetTestBlockTestsArchive returns the bytes of the `.zip` archive containing the tests of a test block
func (useCases *BlocksUseCases) GetTestBlockTestsArchive(dto *dtos.GetBlockTestsArchiveDTO) (archive []byte, err error) {
	// Validate the teacher is the owner of the block
	ownsBlock, err := useCases.BlocksRepository.DoesTeacherOwnsTestBlock(dto.TeacherUUID, dto.BlockUUID)
	if err != nil {
		return nil, err
	}
	if !ownsBlock {
		return nil, blocksErrors.TeacherDoesNotOwnBlock{}
	}

	// Get the UUID of the block's tests archive
	uuid, err := useCases.BlocksRepository.GetTestArchiveUUIDFromTestBlockUUID(dto.BlockUUID)
	if err != nil {
		return nil, err
	}

	// Get the archive from the microservice
	return useCases.StaticFilesRepository.GetArchiveBytes(&staticFilesDTOs.StaticFileArchiveDTO{
		FileUUID: uuid,
		FileType: "test",
	})
}
