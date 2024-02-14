package awsutils

import (
	"context"
	"mime/multipart"

	"gitlab.com/tuneverse/toolkit/core/awsmanager"
)

// CloudService represents a service for managing cloud-related operations.
type CloudService struct {
	awsConf *awsmanager.AwsConfig // awsConf is a configuration for AWS services.
}

// NewCloudService creates a new instance of CloudService with the provided AWS configuration.
func NewCloudService(awsConfig *awsmanager.AwsConfig) CloudServiceImply {
	return &CloudService{
		awsConf: awsConfig,
	}
}

// CloudServiceImply is an interface defining the methods for working with cloud services.
type CloudServiceImply interface {
	DeleteObject(ctx context.Context, bucket, key string) error
	GetObject(ctx context.Context, bucket, key string, expiration int) (string, error)
	UploadToS3(bucketName, key string, fileHeader *multipart.FileHeader, contentType string) error
}
