package services

import (
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-go-std/s3"
)

type StorageService struct {
	s3Service *s3.S3Service
}

func NewStorageService(endpoint string, accessKey string, secretKey string, bucket string, secure bool) (*StorageService, error) {
	s3Service, err := s3.NewS3Service(endpoint, accessKey, secretKey, bucket, secure)
	if err != nil {
		return nil, err
	}

	err = s3Service.MakeBucketIfNotExists(bucket)
	if err != nil {
		return nil, err
	}

	return &StorageService{
		s3Service: s3Service,
	}, nil
}

func (s *StorageService) GetObject(path string) (*minio.Object, error) {
	return s.s3Service.GetObject(path)
}

func (s *StorageService) GetObjects(path string) <-chan *s3.GetObjectsResult {
	return s.s3Service.GetObjects(path)
}

func (s *StorageService) PutSubmissionObject(submissionId string, path string, reader io.Reader, size int64) (info minio.UploadInfo, err error) {
	return s.s3Service.PutObject(fmt.Sprintf("transpilations/%s/%s", submissionId, path), reader, size)
}

func (s *StorageService) ListSubmissionFiles(submissionID string) <-chan *s3.GetObjectsResult {
	ch := make(chan *s3.GetObjectsResult)

	go func() {
		defer close(ch)
		for _, path := range []string{"transpilations/", "transpilations-results/"} {
			objects := s.s3Service.GetObjects(path + submissionID)

			for object := range objects {
				ch <- object
			}
		}
	}()

	return ch
}

func (s *StorageService) DeleteSubmission(id string) error {
	logrus.WithField("id", id).Debug("Deleting submission from S3")
	for _, path := range []string{"transpilations/", "transpilations-results/"} {
		err := s.s3Service.RemoveObjects(path + id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StorageService) SizeofObjects(prefix string) int64 {
	size := int64(0)

	for object := range s.s3Service.GetObjects(prefix) {
		size += object.Size
	}

	return size
}

// Set a tag for the objects to be deleted by Lifecycle later on
func (s *StorageService) ScheduleForDeletion(id string) error {
	for _, path := range []string{"transpilations/", "transpilations-results/"} {
		err := s.s3Service.ScheduleForDeletion(path + id)
		if err != nil {
			return err
		}
	}

	return nil
}
