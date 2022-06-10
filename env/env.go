package env

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Env struct {
	DatabaseDriver   string `env:"DATABASE_DRIVER" env-required:"true"`
	DatabaseEndpoint string `env:"DATABASE_ENDPOINT" env-required:"true"`

	S3Bucket          string `env:"S3_BUCKET" env-required:"true"`
	S3AccessKey       string `env:"S3_ACCESS_KEY" env-required:"true"`
	S3SecretKey       string `env:"S3_SECRET_KEY" env-required:"true"`
	S3Endpoint        string `env:"S3_ENDPOINT" env-required:"true"`
	S3HTTPSEnabled    bool   `env:"S3_HTTPS_ENABLED" env-default:"false"`
	SubmissionsFolder string `env:"SUBMISSIONS_FOLDER" env-required:"true"`

	NSQEndpoint        string `env:"NSQ_ENDPOINT" env-required:"true"`
	NSQLookupdEndpoint string `env:"NSQ_LOOKUPD" env-required:"true"`

	GithubOAuthClientId     string `env:"GITHUB_OAUTH_CLIENT_ID" env-required:"true"`
	GithubOAuthClientSecret string `env:"GITHUB_OAUTH_CLIENT_SECRET" env-required:"true"`

	StripeSecretKey             string `env:"STRIPE_SECRET_KEY" env-required:"true"`
	StripeTierProBase           string `env:"STRIPE_TIER_PRO_BASE" env-required:"true"`
	StripeTierProMetered        string `env:"STRIPE_TIER_PRO_METERED" env-required:"true"`
	StripeTierEnterpriseBase    string `env:"STRIPE_TIER_ENTERPRISE_BASE" env-required:"true"`
	StripeTierEnterpriseMetered string `env:"STRIPE_TIER_ENTERPRISE_METERED" env-required:"true"`
	StripeWebhookSecret         string `env:"STRIPE_WEBHOOK_SECRET" env-required:"true"`

	LogFormat string `env:"LOG_FORMAT" env-default:"json"`
	LogLevel  string `env:"LOG_LEVEL" env-default:"info"`
	SentryDSN string `env:"SENTRY_DSN"`
	Env       string `env:"ENV" env-required:"true"`
}

var env Env

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Warn("Failed to load env variables from file")
	}

	return cleanenv.ReadEnv(&env)
}

func Get() *Env {
	return &env
}
