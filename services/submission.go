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

type TranspilerDetails struct {
	FileExtension string
	Targets       map[string]*TranspilerDetailsTarget
}

type TranspilerDetailsTarget struct {
	FileExtension string
	QueueName     string
}

type SubmissionService struct {
	queueService    *queue.QueueService
	databaseService *DatabaseService
	storageService  *StorageService

	submissionQueues map[string]*TranspilerDetails
}

func NewSubmissionService(queueService *queue.QueueService, databaseService *DatabaseService, storageService *StorageService) *SubmissionService {
	return &SubmissionService{
		queueService:    queueService,
		databaseService: databaseService,
		storageService:  storageService,
		submissionQueues: map[string]*TranspilerDetails{
			"c": {
				FileExtension: ".c",
				Targets: map[string]*TranspilerDetailsTarget{
					"go": {
						FileExtension: ".go",
					},
				},
			},
			"lua": {
				FileExtension: ".lua",
				Targets: map[string]*TranspilerDetailsTarget{
					"ruby": {
						FileExtension: ".rb",
					},
				},
			},
		},
	}
}

type LanguagePairDetails struct {
	SourceLanguageFileExtension string
	TargetLanguageFileExtension string
}

func (s *SubmissionService) GetLanguagePairDetails(sourceLanguage string, targetLanguage string) (*LanguagePairDetails, error) {
	source, ok := s.submissionQueues[sourceLanguage]
	if !ok {
		return nil, fmt.Errorf("source language %s is not supported", sourceLanguage)
	}

	target, ok := source.Targets[targetLanguage]
	if !ok {
		return nil, fmt.Errorf("target language %s is not supported for source language %s", targetLanguage, sourceLanguage)
	}

	return &LanguagePairDetails{
		SourceLanguageFileExtension: source.FileExtension,
		TargetLanguageFileExtension: target.FileExtension,
	}, nil
}

type SubmissionMessage struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

type SubmissionStatusMessage struct {
	ID        string            `json:"id"`
	Status    submission.Status `json:"status"`
	Reason    string            `json:"reason"`
	Timestamp int64             `json:"timestamp"`
}

func (s *SubmissionService) PublishSubmissionToTranspile(sub SubmissionMessage) error {
	bytes, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("transpilation_jobs_%s_to_%s", sub.SourceLanguage, sub.TargetLanguage)
	return s.queueService.Publish(topic, bytes)
}

func (s *SubmissionService) HandleSubmissionStatus(msg SubmissionStatusMessage) error {
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

	receivedAt := time.UnixMilli(msg.Timestamp)

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
