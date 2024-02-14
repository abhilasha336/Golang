package awsutils

import (
	"context"

	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/tuneverse/toolkit/core/logger"
)

func (cloud *CloudService) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := cloud.awsConf.S3().DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("DeleteObject failed, err=%s", err.Error())
		return err
	}
	return nil
}

// GetObject retrieves a pre-signed URL for an object in an S3 bucket.
func (cloud *CloudService) GetObject(ctx context.Context, bucket, key string, expiration int) (string, error) {

	req, err := s3.NewPresignClient(cloud.awsConf.S3()).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiration * int(time.Hour))
	})

	//Generate pre-signed URL
	if err != nil {
		logger.Log().WithContext(ctx).Errorf("GetObject failed, err=%s", err.Error())
		return "", err
	}

	return req.URL, nil

}

// UploadToS3 uploads a file to an AWS S3 bucket with the specified key and content type.
func (cl *CloudService) UploadToS3(bucketName, key string, fileHeader *multipart.FileHeader, contentType string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err // Return the error if unable to open the file
	}
	defer file.Close()

	_, err = cl.awsConf.Uploader().Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return err // Return the error
	}
	return nil
}
