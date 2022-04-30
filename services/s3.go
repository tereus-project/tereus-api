package services

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tereus-project/tereus-api/env"
)

type S3Service struct {
	client *minio.Client
}

func NewS3Service(endpoint string, accessKey string, secretKey string) (*S3Service, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &S3Service{
		client: client,
	}, nil
}

func (s *S3Service) MakeBucketIfNotExists(name string) error {
	exists, err := s.client.BucketExists(context.Background(), name)
	if err != nil {
		log.Fatalln(err)
	}

	if !exists {
		err = s.client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			return err
		}
	}

	return nil
}

func (s *S3Service) PutObject(bucket string, path string, reader io.Reader, size int64) (info minio.UploadInfo, err error) {
	return s.client.PutObject(
		context.Background(),
		env.S3Bucket,
		path,
		reader,
		size,
		minio.PutObjectOptions{},
	)
}

func (s *S3Service) GetObject(bucket string, path string) (*minio.Object, error) {
	return s.client.GetObject(context.Background(), bucket, path, minio.GetObjectOptions{})
}

type GetObjectsResult struct {
	Err  error
	Path string
}

func (s *S3Service) GetObjects(bucket string, prefix string) (paths <-chan *GetObjectsResult, err error) {
	ch := make(chan *GetObjectsResult)

	go func() {
		defer close(ch)

		for object := range s.client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
			ch <- &GetObjectsResult{
				Err:  object.Err,
				Path: object.Key,
			}
		}
	}()

	return ch, nil
}
