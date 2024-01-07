package implementations

import (
	"log"

	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
)

type UpdatesChannel chan *dtos.SubmissionStatusUpdateDTO

type SubmissionsRealTimeUpdatesSender struct {
	// New updates are pushed to this channel
	Updates UpdatesChannel

	// All the open connections
	OpenConnections map[string]UpdatesChannel
}

// Singleton instance
var submissionsRealTimeUpdatesSenderInstance *SubmissionsRealTimeUpdatesSender

func GetSubmissionsRealTimeUpdatesSenderInstance() *SubmissionsRealTimeUpdatesSender {
	if submissionsRealTimeUpdatesSenderInstance == nil {
		submissionsRealTimeUpdatesSenderInstance = &SubmissionsRealTimeUpdatesSender{
			Updates:         make(UpdatesChannel),
			OpenConnections: make(map[string]UpdatesChannel),
		}
	}

	return submissionsRealTimeUpdatesSenderInstance
}

// Create a channel for a submission
func (sender *SubmissionsRealTimeUpdatesSender) CreateChannel(submissionUUID string) *UpdatesChannel {
	currentChannel, channelExists := sender.OpenConnections[submissionUUID]
	if channelExists {
		return &currentChannel
	}

	newChannel := make(UpdatesChannel)
	sender.OpenConnections[submissionUUID] = newChannel

	log.Printf(
		"Created a new channel for submission %s. Total channels: %d",
		submissionUUID,
		len(sender.OpenConnections),
	)
	return &newChannel
}

// Delete a channel for a submission
func (sender *SubmissionsRealTimeUpdatesSender) DeleteChannel(submissionUUID string) {
	_, channelExists := sender.OpenConnections[submissionUUID]
	if !channelExists {
		return
	}

	close(sender.OpenConnections[submissionUUID])
	delete(sender.OpenConnections, submissionUUID)

	log.Printf(
		"Deleted channel for submission %s. Total channels: %d",
		submissionUUID,
		len(sender.OpenConnections),
	)
}

// SendUpdate sends an update to the updates channel
func (sender *SubmissionsRealTimeUpdatesSender) SendUpdate(update *dtos.SubmissionStatusUpdateDTO) {
	sender.Updates <- update
}

// listen listens for new updates and new connections
func (sender *SubmissionsRealTimeUpdatesSender) Listen() {
	log.Println("[SSE]: Listening for new updates")

	for update := range sender.Updates {
		log.Printf(
			"[SSE]: Received update for submission %s: %s",
			update.SubmissionUUID,
			update.SubmissionStatus,
		)

		// Send the update to the corresponding channel
		ch, chExists := sender.OpenConnections[update.SubmissionUUID]
		if chExists {
			ch <- update
			log.Printf(
				"Sent update to submission %s",
				update.SubmissionUUID,
			)
		}
	}
}
