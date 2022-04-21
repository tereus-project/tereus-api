package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	DatabaseDriver          string
	DatabaseEndpoint        string
	S3Bucket                string
	S3AccessKey             string
	S3SecretKey             string
	S3Endpoint              string
	RabbitMQEndpoint        string
	SubmissionsFolder       string
	GithubOAuthClientId     string
	GithubOAuthClientSecret string
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	DatabaseDriver = os.Getenv("DATABASE_DRIVER")
	if DatabaseDriver == "" {
		return fmt.Errorf("DATABASE_DRIVER is not set")
	}

	DatabaseEndpoint = os.Getenv("DATABASE_ENDPOINT")
	if DatabaseEndpoint == "" {
		return fmt.Errorf("DATABASE_ENDPOINT is not set")
	}

	S3Bucket = os.Getenv("S3_BUCKET")
	if S3Bucket == "" {
		return fmt.Errorf("S3_BUCKET is not set")
	}

	S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	if S3AccessKey == "" {
		return fmt.Errorf("S3_ACCESS_KEY is not set")
	}

	S3SecretKey = os.Getenv("S3_SECRET_KEY")
	if S3SecretKey == "" {
		return fmt.Errorf("S3_SECRET_KEY is not set")
	}

	S3Endpoint = os.Getenv("S3_ENDPOINT")
	if S3Endpoint == "" {
		return fmt.Errorf("S3_ENDPOINT is not set")
	}

	RabbitMQEndpoint = os.Getenv("RABBITMQ_ENDPOINT")
	if RabbitMQEndpoint == "" {
		return fmt.Errorf("RABBITMQ_ENDPOINT is not set")
	}

	SubmissionsFolder = os.Getenv("SUBMISSIONS_FOLDER")
	if SubmissionsFolder == "" {
		return fmt.Errorf("SUBMISSIONS_FOLDER is not set")
	}

	GithubOAuthClientId = os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	if GithubOAuthClientId == "" {
		return fmt.Errorf("GITHUB_OAUTH_CLIENT_ID is not set")
	}

	GithubOAuthClientSecret = os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
	if GithubOAuthClientSecret == "" {
		return fmt.Errorf("GITHUB_OAUTH_CLIENT_SECRET is not set")
	}

	return nil
}
