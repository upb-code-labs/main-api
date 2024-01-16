package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
)

type LaboratoriesRepository interface {
	GetLaboratoryByUUID(uuid string) (laboratory *entities.Laboratory, err error)
	GetLaboratoryInformationByUUID(uuid string) (laboratory *dtos.LaboratoryDetailsDTO, err error)
	SaveLaboratory(dto *dtos.CreateLaboratoryDTO) (laboratory *entities.Laboratory, err error)
	UpdateLaboratory(dto *dtos.UpdateLaboratoryDTO) error

	CreateMarkdownBlock(laboratoryUUID string) (blockUUID string, err error)
	CreateTestBlock(dto *dtos.CreateTestBlockDTO) (blockUUID string, err error)

	GetTotalTestBlocks(laboratoryUUID string) (total int, err error)
	GetStudentsProgress(laboratoryUUID string) (progress []*dtos.LaboratoryStudentProgressDTO, err error)

	DoesTeacherOwnLaboratory(teacherUUID string, laboratoryUUID string) (bool, error)
}
