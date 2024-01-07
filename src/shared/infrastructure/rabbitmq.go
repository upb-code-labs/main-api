package infrastructure

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitMQChannel *amqp.Channel
var rabbitMQSubmissionsQueue *amqp.Queue

func ConnectToRabbitMQ() {
	// Stablish connection
	rabbitMQConnectionString := GetEnvironment().RabbitMQConnectionString
	conn, err := amqp.Dial(rabbitMQConnectionString)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set channel
	log.Println("Connected to RabbitMQ")
	rabbitMQChannel = ch
}

func CloseRabbitMQConnection() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
}

func GetRabbitMQChannel() *amqp.Channel {
	if rabbitMQChannel == nil {
		ConnectToRabbitMQ()
	}

	return rabbitMQChannel
}

func GetRabbitMQSubmissionsQueue() *amqp.Queue {
	if rabbitMQSubmissionsQueue == nil {
		ch := GetRabbitMQChannel()

		// Declare queue
		qName := "submissions"
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
			log.Fatal(err.Error())
		}

		// Set fair dispatch
		maxPrefetchCount := 4 // Limit to 4 submissions per worker
		err = ch.Qos(
			maxPrefetchCount,
			0,
			false,
		)

		if err != nil {
			log.Fatal(err.Error())
		}

		// Set queue
		log.Println("RabbitMQ submissions queue declared / set")
		rabbitMQSubmissionsQueue = &q
	}

	return rabbitMQSubmissionsQueue
}
