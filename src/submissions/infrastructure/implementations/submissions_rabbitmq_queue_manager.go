package implementations

import (
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SubmissionsRabbitMQQueueManager struct {
	SubmissionsQueue *amqp.Queue
}

// Singleton
var submissionsRabbitMQQueueManagerInstance *SubmissionsRabbitMQQueueManager

func GetSubmissionsRabbitMQQueueManagerInstance() *SubmissionsRabbitMQQueueManager {
	if submissionsRabbitMQQueueManagerInstance == nil {
		submissionsRabbitMQQueueManagerInstance = &SubmissionsRabbitMQQueueManager{
			SubmissionsQueue: sharedInfrastructure.GetRabbitMQSubmissionsQueue(),
		}
	}

	return submissionsRabbitMQQueueManagerInstance
}

// Methods implementation
func (queueManager *SubmissionsRabbitMQQueueManager) QueueWork(work *entities.SubmissionWork) (err error) {
	return nil
}
