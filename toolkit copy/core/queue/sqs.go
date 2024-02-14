package queue

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gitlab.com/tuneverse/toolkit/core/awsmanager"
)

// SQSConfig contains configurations for creating a queue and receiving a message from SQS.
type SQSConfig struct {

	// Configuration for creating a queue
	QueueInfo *sqs.CreateQueueInput

	// Configuration for receiving a message from SQS
	ReceiveMessageConfig *sqs.ReceiveMessageInput
}

type sqsQueue struct {
	client               *sqs.Client
	queueInfo            *sqs.CreateQueueOutput
	receiveMessageConfig *sqs.ReceiveMessageInput
}

// NewSQSQueue creates a new SQS queue with the given configurations.
func NewSQSQueue(awsConfig *awsmanager.AwsConfig, config *SQSConfig, optFns ...func(*sqs.Options)) (Queue, error) {
	client := awsConfig.SQS(optFns...)
	queueData, err := create(client, config.QueueInfo)
	if err != nil {
		return nil, err
	}
	return &sqsQueue{
		client:               client,
		queueInfo:            queueData,
		receiveMessageConfig: config.ReceiveMessageConfig,
	}, nil
}

// Compose a message to the queue.
func (rsqsSvc *sqsQueue) ComposeMessage(ctx context.Context, messageBody []byte) (interface{}, error) {
	var message interface{}
	messageBodyString := string(messageBody)
	message = &sqs.SendMessageInput{
		MessageBody:  aws.String(messageBodyString),
		DelaySeconds: *aws.Int32(10),
	}
	return message, nil
}

// create creates a new queue with the given configurations.
func create(client *sqs.Client, config *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	if config.QueueName == nil {
		return nil, errors.New("queue name is required")
	}
	sqsOutput, err := client.CreateQueue(context.TODO(), config)
	if err != nil {
		return nil, err
	}

	// If the CreateQueue operation is successful, sqsOutput.QueueUrl should be populated.
	if sqsOutput.QueueUrl == nil || *sqsOutput.QueueUrl == "" {
		return nil, errors.New("failed to create the queue. Queue URL is empty")
	}

	return sqsOutput, nil
}

// Delete deletes a message from an Amazon SQS queue.
func (sqsSvc *sqsQueue) Delete(ctx context.Context, receiptHandle string) error {

	if receiptHandle == "" {
		return errors.New("receipt handle is required")
	}

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      sqsSvc.queueInfo.QueueUrl,
		ReceiptHandle: &receiptHandle,
	}

	_, err := sqsSvc.client.DeleteMessage(context.TODO(), dMInput)
	if err != nil {
		return err
	}

	return nil
}

// Send sends a message to the queue.
func (sqsSvc *sqsQueue) Send(ctx context.Context, input interface{}) error {

	data, ok := input.(*sqs.SendMessageInput)
	if !ok {
		return errors.New("invalid input")
	}

	data.QueueUrl = sqsSvc.queueInfo.QueueUrl
	_, err := sqsSvc.client.SendMessage(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// Receive receives a message from the queue.
func (sqsSvc *sqsQueue) Receive(ctx context.Context) (interface{}, error) {
	sqsSvc.receiveMessageConfig.QueueUrl = sqsSvc.queueInfo.QueueUrl
	return sqsSvc.client.ReceiveMessage(ctx, sqsSvc.receiveMessageConfig)
}

// Close closes the connection to the SQS server.
func (sqsSvc *sqsQueue) Close() error {
	return nil
}
