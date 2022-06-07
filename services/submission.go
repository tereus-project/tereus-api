package services

import (
	"fmt"

	"github.com/tereus-project/tereus-api/ent/submission"
)

type SubmissionService struct {
	queueService *QueueService

	submissionQueues map[string]map[string]*Queue[SubmissionMessage]
}

func NewSubmissionService(queueService *QueueService) *SubmissionService {
	return &SubmissionService{
		queueService: queueService,
		submissionQueues: map[string]map[string]*Queue[SubmissionMessage]{
			"c": {
				"go": NewQueue[SubmissionMessage]("remix_jobs_c_to_go", queueService),
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

func (s *SubmissionService) PublishSubmission(message *SubmissionMessage) error {
	return s.submissionQueues[message.SourceLanguage][message.TargetLanguage].Publish(message.ID, message)
}

type SubmissionStatusMessage struct {
	ID     string            `json:"id"`
	Status submission.Status `json:"status"`
	Reason string            `json:"reason"`
}

func (s *SubmissionService) ConsumeSubmissionsStatus() <-chan SubmissionStatusMessage {
	queue := NewQueue[SubmissionStatusMessage]("submission_status", s.queueService)
	return queue.Consume()
}
