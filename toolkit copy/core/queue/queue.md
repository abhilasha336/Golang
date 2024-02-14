# Package Queue
This package provides a generic interface and implementations for working with various queue service providers. The supported providers currently include AWS Simple Queue Service (SQS) and RabbitMQ.

# Usage

```go
type Queue interface {
    ComposeMessage(ctx context.Context, message []byte) (interface{}, error)
	Send(ctx context.Context, message interface{}) error
	Receive(ctx context.Context) (interface{}, error)
	Delete(ctx context.Context, receiptHandle string) error
	Close() error
}
```
The `Queue` interface defines methods for interacting with queues.

# RabbitMQ

### Configuration

To use `RabbitMQ`, create an instance of `RabbitMQConfig` with the necessary configurations, such as the connection URL, queue name, and other optional settings.

**Creating a RabbitMQ Queue**

```go
cfg := &queue.RabbitMQConfig{
    // ... (provide RabbitMQ configurations)
}

mq, err := queue.NewRabbitMQQueue(cfg)
if err != nil {
    log.Fatal(err)
}
defer mq.Close()
```

***Composing a Message***
```go
messageBody := []byte("your_message_data_here")
message, err := mq.ComposeMessage(context.Background(), messageBody)
if err != nil {
    log.Fatal(err)
}
```

***Sending a Message***

```go
message := amqp091.Publishing{
    // ... (provide message details)
}

err := mq.Send(context.Background(), message)
if err != nil {
    log.Fatal(err)
}
```

***Receiving a Message***

```go
msg, err := mq.Receive(context.Background())
if err != nil {
    log.Fatal(err)
}

//when receiving message make sure to convert to appropriate types, otherwise it will return an error

data, ok := msg.(<-chan amqp091.Delivery)
if !ok {
    //handle it accordingly
}
//otherwise process the received messages according to usecases.
```
**Process the received message**

***Deleting a Message***

```go
err := mq.Delete(context.Background(), receiptHandle)
if err != nil {
    log.Fatal(err)
}
```

# AWS SQS

### Configuration

To use AWS SQS, create an instance of `SQSConfig` with the necessary configurations, such as `CreateQueueInput` for creating a queue and `ReceiveMessageInput` for receiving a message.

***Creating an SQS Queue***

```go
awsConfig := &awsmanager.AwsConfig{
    // ... (provide AWS configurations)
}

sqsConfig := &queue.SQSConfig{
    QueueInfo: &sqs.CreateQueueInput{
        // ... (provide queue creation details)
    },
    ReceiveMessageConfig: &sqs.ReceiveMessageInput{
        // ... (provide receive message details)
    },
}

sqs, err := queue.NewSQSQueue(awsConfig, sqsConfig)
if err != nil {
    log.Fatal(err)
}
```
***Sending a Message***

```go
message := &sqs.SendMessageInput{
    // ... (provide message details)
}

err := sqs.Send(context.Background(), message)
if err != nil {
    log.Fatal(err)
}
```

***Receiving a Message***

```go
msg, err := sqs.Receive(context.Background())
if err != nil {
    log.Fatal(err)
}

//when receiving message make sure to convert to appropriate types, otherwise it will return an error

data, ok := msg.(*sqs.ReceiveMessageOutput)
if !ok {
    //handle it accordingly
}
```

**Process the received message**

***Deleting a Message***

After processing the message make sure to remove it from the queue.

```go
err := sqs.Delete(context.Background(), receiptHandle)
if err != nil {
    log.Fatal(err)
}
//Make sure to handle errors appropriately in your application logic.
```