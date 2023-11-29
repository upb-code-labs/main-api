package application

import (
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/errors"
)

type BlocksUseCases struct {
	BlocksRepository definitions.BlockRepository
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
