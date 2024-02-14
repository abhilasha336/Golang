package awsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AwsConfig struct {
	aws *aws.Config
}

func CreateAwsSession(optFns ...func(*config.LoadOptions) error) (*AwsConfig, error) {
	awsConfig := &AwsConfig{}
	defaultConf, err := config.LoadDefaultConfig(context.TODO(), optFns...)
	if err != nil {
		return nil, err
	}
	conf := aws.Config(defaultConf)
	if !isAccessCredentialExist(&conf) {
		return nil, fmt.Errorf("unable to load aws credentials")
	}
	awsConfig.aws = &conf
	return awsConfig, nil
}

func (conf AwsConfig) S3(optFns ...func(*s3.Options)) *s3.Client {
	return s3.NewFromConfig(*conf.aws, optFns...)
}

func (conf AwsConfig) Uploader(options ...func(*manager.Uploader)) *manager.Uploader {
	return manager.NewUploader(conf.S3(), options...)
}

func (conf AwsConfig) SQS(optFns ...func(*sqs.Options)) *sqs.Client {
	return sqs.NewFromConfig(*conf.aws, optFns...)
}

func isAccessCredentialExist(conf *aws.Config) bool {
	cred, err := conf.Credentials.Retrieve(context.Background())
	if err == nil {
		return len(cred.AccessKeyID) != 0 && len(cred.SecretAccessKey) != 0
	}
	return false
}

func WithCredentialsProvider(accessKey, accessSecret string) config.LoadOptionsFunc {
	return config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, accessSecret, ""))
}

func WithRegion(region string) config.LoadOptionsFunc {
	return config.WithRegion(region)
}
