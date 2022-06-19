package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-go-std/queue"
)

type SubmissionService struct {
	queueService    *queue.QueueService
	databaseService *DatabaseService
	storageService  *StorageService

	submissionQueues map[string]map[string]string
}

func NewSubmissionService(queueService *queue.QueueService, databaseService *DatabaseService, storageService *StorageService) *SubmissionService {
	return &SubmissionService{
		queueService:    queueService,
		databaseService: databaseService,
		storageService:  storageService,
		submissionQueues: map[string]map[string]string{
			"c": {
				"go": "transpilation_jobs_c_to_go",
			},
		},
	}
}

func (s *SubmissionService) CheckSupport(sourceLanguage string, targetLanguage string) error {
	m, ok := s.submissionQueues[sourceLanguage]
	if !ok {
		return fmt.Errorf("source language %s is not supported", sourceLanguage)
	}

	if _, ok := m[targetLanguage]; !ok {
		return fmt.Errorf("target language %s is not supported for source language %s", targetLanguage, sourceLanguage)
	}

	return nil
}

type SubmissionMessage struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

type SubmissionStatusMessage struct {
	ID     string            `json:"id"`
	Status submission.Status `json:"status"`
	Reason string            `json:"reason"`
}

func (s *SubmissionService) PublishSubmissionToTranspile(sub SubmissionMessage) error {
	bytes, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	return s.queueService.Publish(s.submissionQueues[sub.SourceLanguage][sub.TargetLanguage], bytes)
}

func (s *SubmissionService) HandleSubmissionStatus(msg SubmissionStatusMessage, receivedAt time.Time) error {
	logrus.WithField("status", msg).Info("Handling submission status")

	id, err := uuid.Parse(msg.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse submission ID")
		return nil
	}

	var submissionBytesCount int64
	if msg.Status == submission.StatusDone {
		submissionBytesCount = s.storageService.SizeofObjects(fmt.Sprintf("transpilations-results/%s/", id))
	}

	submissionUpdate := s.databaseService.Submission.
		Update().
		Where(
			submission.ID(id),
			submission.StatusIn(submission.StatusPending, submission.StatusProcessing),
		).
		SetSubmissionTargetSizeBytes(int(submissionBytesCount)).
		SetStatus(msg.Status).
		SetReason(msg.Reason)

	switch msg.Status {
	case submission.StatusProcessing:
		submissionUpdate = submissionUpdate.SetProcessingStartedAt(receivedAt)
	case submission.StatusDone, submission.StatusFailed:
		submissionUpdate = submissionUpdate.SetProcessingFinishedAt(receivedAt)
	}

	err = submissionUpdate.Exec(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to update submission status")
		return err
	}

	return nil
}
