package implementations

import (
	"context"
	"time"

	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/errors"
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	channel := sharedInfrastructure.GetRabbitMQChannel()

	// Parse work to JSON
	stringifiedWork, err := work.ToJSON()
	if err != nil {
		return err
	}

	// Publish work to queue
	msgExchange := ""
	msgKey := queueManager.SubmissionsQueue.Name
	msgMandatory := false
	msgImmediate := false

	err = channel.PublishWithContext(
		ctx,
		msgExchange,
		msgKey,
		msgMandatory,
		msgImmediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(stringifiedWork),
		},
	)

	if err != nil {
		return errors.UnableToQueueSubmissionWork{}
	}

	return nil
}
