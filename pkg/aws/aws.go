package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type LocalAwsCredentialProvider struct {
}

func (me *LocalAwsCredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		SessionToken:    "test",
		CanExpire:       false,
		Source:          "test",
	}, nil
}

func V2Config() aws.Config {

	cfg := aws.NewConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewCredentialsCache(&LocalAwsCredentialProvider{})

	config.WithCredentialsCacheOptions(func(o *aws.CredentialsCacheOptions) {
		o.ExpiryWindow = 0

	})

	return *cfg
}
