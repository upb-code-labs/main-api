package implementations

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type BlocksPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var blocksPostgresRepositoryInstance *BlocksPostgresRepository

func GetBlocksPostgresRepositoryInstance() *BlocksPostgresRepository {
	if blocksPostgresRepositoryInstance == nil {
		blocksPostgresRepositoryInstance = &BlocksPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
		}
	}

	return blocksPostgresRepositoryInstance
}

func (repository *BlocksPostgresRepository) UpdateMarkdownBlockContent(blockUUID string, content string) (err error) {
	return nil
}
