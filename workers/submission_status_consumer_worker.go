package workers

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/services"
	"github.com/tereus-project/tereus-go-std/queue"
)

type SubsmissionStatusHandler struct {
	submissionService *services.SubmissionService
}

// HandleMessage implements the Handler interface.
// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
func (h *SubsmissionStatusHandler) HandleMessage(m *nsq.Message) error {
	logrus.WithField("nsq_msg_id", m.ID).Info("Received submission status message")

	var msg services.SubmissionStatusMessage
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		logrus.WithError(err).Error("Error unmarshaling message")
		return nil
	}

	return h.submissionService.HandleSubmissionStatus(msg)
}

func RegisterStatusConsumerWorker(submissionService *services.SubmissionService, queueService *queue.QueueService) error {
	logrus.Info("Starting submission status consumer worker")

	h := &SubsmissionStatusHandler{
		submissionService: submissionService,
	}

	err := queueService.AddHandler("transpilation_submission_status", "api", h.HandleMessage)

	return err
}
