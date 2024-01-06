package implementations

import (
	"encoding/json"
	"log"

	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SubmissionsRealTimeUpdatesQueueMgr struct {
	Queue   *amqp.Queue
	Channel <-chan amqp.Delivery
}

// ## Singleton instance ##

var submissionsRealTimeUpdatesQueueMgrInstance *SubmissionsRealTimeUpdatesQueueMgr

func GetSubmissionsRealTimeUpdatesQueueMgrInstance() *SubmissionsRealTimeUpdatesQueueMgr {
	if submissionsRealTimeUpdatesQueueMgrInstance == nil {
		submissionsRealTimeUpdatesQueueMgrInstance = &SubmissionsRealTimeUpdatesQueueMgr{
			Queue: getSubmissionsRealTimeUpdatesQueue(),
			// Channel will be set when ListenForUpdates is called
		}
	}

	return submissionsRealTimeUpdatesQueueMgrInstance
}

// ## Public methods ##

// ListenForUpdates listens for updates in the submission real time updates queue
func (queueMgr *SubmissionsRealTimeUpdatesQueueMgr) ListenForUpdates() {
	ch := sharedInfrastructure.GetRabbitMQChannel()

	// Consume messages
	qName := queueMgr.Queue.Name
	qConsumer := ""
	qAutoAck := false
	qExclusive := false
	qNoLocal := false
	qNoWait := false
	qArgs := amqp.Table{}

	msgs, err := ch.Consume(
		qName,
		qConsumer,
		qAutoAck,
		qExclusive,
		qNoLocal,
		qNoWait,
		qArgs,
	)

	if err != nil {
		log.Fatal(
			"[RabbitMQ]: There was an error while consuming messages from the submission real time updates queue",
			err.Error(),
		)
	}

	// Set the channel
	queueMgr.Channel = msgs

	// Process
	queueMgr.processUpdates()
}

// ## Private methods ##

// getSubmissionsRealTimeUpdatesQueue returns the submissions real time updates queue
func getSubmissionsRealTimeUpdatesQueue() *amqp.Queue {
	noQueueHasBeenDeclared := submissionsRealTimeUpdatesQueueMgrInstance == nil ||
		submissionsRealTimeUpdatesQueueMgrInstance.Queue == nil

	if noQueueHasBeenDeclared {
		ch := sharedInfrastructure.GetRabbitMQChannel()

		// Declare queue
		qName := "submission-real-time-updates"
		qDurable := true
		qAutoDelete := false
		qExclusive := false
		qNoWait := false
		qArgs := amqp.Table{}

		q, err := ch.QueueDeclare(
			qName,
			qDurable,
			qAutoDelete,
			qExclusive,
			qNoWait,
			qArgs,
		)

		if err != nil {
			log.Fatal(
				"[RabbitMQ]: There was an error while declaring the submission real time updates queue",
				err.Error(),
			)
		}

		// Set fair dispatch
		maxPrefetchCount := 4 // Limit to 4 updates per worker
		err = ch.Qos(
			maxPrefetchCount,
			0,
			false,
		)
		if err != nil {
			log.Fatal(
				"[RabbitMQ]: There was an error while setting the fair dispatch for the submission real time updates queue",
				err.Error(),
			)
		}

		log.Println("[RabbitMQ]: Submission real time updates queue declared")
		return &q
	}

	return submissionsRealTimeUpdatesQueueMgrInstance.Queue
}

// processUpdates processes the updates received from the submission real time updates queue
func (queueMgr *SubmissionsRealTimeUpdatesQueueMgr) processUpdates() {
	log.Println("[RabbitMQ Submissions Real Time Updates Queue]: Listening for updates...")

	for msg := range queueMgr.Channel {
		go queueMgr.processUpdate(msg)
	}
}

// processUpdate processes a single update received from the submission real time updates queue
func (queueMgr *SubmissionsRealTimeUpdatesQueueMgr) processUpdate(msg amqp.Delivery) {
	// Acknowledge the message after processing it
	defer msg.Ack(false)

	// Get the update
	dto := dtos.SubmissionStatusUpdateDTO{}
	err := json.Unmarshal(msg.Body, &dto)
	if err != nil {
		log.Println(
			"[RabbitMQ Submissions Real Time Updates Queue]: There was an error while parsing the submission real time update",
			err.Error(),
		)
		return
	}

	// Send the update to the real time updates sender
	realTimeUpdater := GetSubmissionsRealTimeUpdatesSenderInstance()
	realTimeUpdater.SendUpdate(&dto)
}
