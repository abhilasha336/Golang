## AWS Manager
This Go package provides utility functions for managing AWS configurations and sessions, allowing easy setup of AWS SDK v2 clients with customizable options.

### Overview
This Go package provides logging functionalities to simplify logging operations within your application. It includes functions for setting log levels, including dump data, and making API requests to log external services.

#### Table of Contents

- Features
- Usage


### Features

- [func CreateAwsSession(optFns ...func(*config.LoadOptions) error) ( *AwsConfig, error)](#CreateAwsConfig)
- [func (conf AwsConfig) S3(optFns ...func(*s3.Options)) *s3.Client](#CreateS3Client)
- [func (conf AwsConfig) Uploader(options ...func(*manager.Uploader)) *manager.Uploader](#CreateS3UploadManager)
- [func (conf AwsConfig) SQS(optFns ...func(*sqs.Options)) *sqs.Client](#CreateS3UploadManager)
- [func WithCredentialsProvider(accessKey, accessSecret string) config.LoadOptionsFunc](#WithCredential)
- [func WithRegion(region string) config.LoadOptionsFunc](#WithCredential)


### Usage 

    // Example usage of the AWS Manager package

    // Initialize AWS session with default configuration
    conf, err := awsmanager.CreateAwsSession()
    if err != nil {
        fmt.Println("Failed to create AWS session:", err)
        return
    }

    // Use the AWS config for AWS SDK v2 clients
    // Example: s3Client := s3.NewFromConfig(*conf)



    //Example for development mode   
    conf,err:= CreateAwsSession(awsmanager.WithCredentialsProvider(cfg.AWS.AccessKey, cfg.AWS.AccessSecret),awsmanager.WithRegion(cfg.AWS.Region))
	if err != nil {
	 	fmt.Println("Failed to create AWS session:", err)
        return
	}
    // Use the AWS config for AWS SDK v2 clients
    // Example: s3Client := s3.NewFromConfig(*conf)





