package queue

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConfig contains various configurations for the RabbitMQ queue
type RabbitMQConfig struct {

	// URL of the connection string
	URL string

	// The name of the queue, which serves as a unique identifier for the queue.
	Name string

	// Durable (queue survives server restart).
	// If set to true, the queue will survive server restarts; otherwise,
	// it is non-durable and will be deleted upon server restart.
	Durable bool

	// AutoDelete (queue is deleted when no longer in use).
	// If set to true, the queue is automatically deleted when it is no longer in use.
	AutoDelete bool

	// Exclusive (queue only accessible by the declaring connection).
	// If set to true, the queue is exclusive to the declaring connection
	// and cannot be accessed by other connections.
	Exclusive bool

	// NoWait (no wait for a server response).
	// If set to true, the method will not wait for a server response, making it non-blocking.
	NoWait bool

	// Internal specifies whether the exchange should be marked as an internal exchange.
	// Internal exchanges are used by RabbitMQ internals and are not meant for
	// direct client interaction.
	Internal bool

	// Additional arguments for the queue,
	// such as message TTL ("x-message-ttl") and other configuration options.
	Arguments amqp.Table

	////***************important********////

	// Configuration for sending a message to the queue

	// Control flag for message publishing. If true, the message must be routed to a queue.
	Mandatory bool

	// Control flag for message publishing. If true, the message should be
	// delivered to consumers immediately or return an error if there are no consumers.
	Immediate bool

	////***************important********////

	// Configuration for receiving a message from the queue

	// - `NoLocal` specifies that the server should not deliver messages
	// published by this consumer.
	NoLocal bool
	// Control flags for exclusive access and message acknowledgment.
	// - `AutoAck` specifies whether messages are automatically acknowledged upon consumption.
	AutoAck bool
	// A unique identifier (consumer tag) for this consumer, allowing tracking and management
	// of multiple consumers on the same channel.
	Consumer string
	// Additional arguments and configuration options for message consumption,
	// such as message headers and other properties.
	Args amqp.Table
}

// RabbitMQQueue contains a reference to the AMQP connection used for communication with the RabbitMQ server,
// a reference to the AMQP channel, and the configurations for the queue.
type RabbitMQQueue struct {

	// Configurations
	config *RabbitMQConfig

	// A reference to the AMQP connection used for communication with the RabbitMQ server.
	connection *amqp.Connection

	// A reference to the AMQP channel, which is a communication channel
	// within the connection for operations like message publishing and consumption.
	channel *amqp.Channel

	// Queue information
	queueInfo *amqp.Queue
}

// NewRabbitMQQueue creates a new RabbitMQ queue with the given configurations.
// It returns a `Queue` interface and an error.
func NewRabbitMQQueue(cfg *RabbitMQConfig) (Queue, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	mq := &RabbitMQQueue{
		connection: conn,
		channel:    ch,
		config:     cfg,
	}

	queueInfo, err := mq.Create()
	if err != nil {
		return nil, fmt.Errorf("cannot initialise Queue %w", err)
	}

	mq.queueInfo = queueInfo

	return mq, nil
}

// Close closes the channel and connection to the RabbitMQ server.
func (rabbitMQQueue *RabbitMQQueue) Close() error {
	if rabbitMQQueue.channel != nil {
		if err := rabbitMQQueue.channel.Close(); err != nil {
			return err
		}
	}

	if rabbitMQQueue.connection != nil {
		if err := rabbitMQQueue.connection.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Create creates a new queue with the given configurations.
// It returns a pointer to the `amqp.Queue` struct and an error.
func (rabbitMQQueue *RabbitMQQueue) Create() (*amqp.Queue, error) {
	// Check if the `config` parameter is nil
	if rabbitMQQueue.config == nil {
		return nil, errors.New("config must be provided")
	}

	if rabbitMQQueue.config.Name == "" {
		return nil, errors.New("queue name must be present")
	}

	q, err := rabbitMQQueue.channel.QueueDeclare(
		rabbitMQQueue.config.Name,       // Queue name
		rabbitMQQueue.config.Durable,    // Durable
		rabbitMQQueue.config.AutoDelete, // AutoDelete
		rabbitMQQueue.config.Exclusive,  // Exclusive
		rabbitMQQueue.config.NoWait,     // NoWait
		rabbitMQQueue.config.Arguments,  // Arguments
	)

	if err != nil {
		return nil, err
	}

	return &q, nil
}

// Compose a message to the queue.
func (rabbitMQQueue *RabbitMQQueue) ComposeMessage(ctx context.Context, messageBody []byte) (interface{}, error) {
	message := amqp.Publishing{
		ContentType: "application/json",
		Priority:    1,
		Timestamp:   time.Now(),
		Body:        messageBody,
	}
	return message, nil
}

// Send sends a message to the queue.
func (rabbitMQQueue *RabbitMQQueue) Send(ctx context.Context, input interface{}) error {
	data, ok := input.(amqp.Publishing)
	if !ok {
		return errors.New("invalid input")
	}

	return rabbitMQQueue.channel.PublishWithContext(
		ctx,
		"",
		rabbitMQQueue.queueInfo.Name,
		rabbitMQQueue.config.Mandatory,
		rabbitMQQueue.config.Immediate,
		data,
	)
}

// Receive receives a message from the queue.
func (rabbitMQQueue *RabbitMQQueue) Receive(ctx context.Context) (interface{}, error) {
	msg, err := rabbitMQQueue.channel.ConsumeWithContext(
		ctx,
		rabbitMQQueue.config.Name,
		rabbitMQQueue.config.Consumer,
		rabbitMQQueue.config.AutoAck,
		rabbitMQQueue.config.Exclusive,
		rabbitMQQueue.config.NoLocal,
		rabbitMQQueue.config.NoWait,
		rabbitMQQueue.config.Args,
	)

	return msg, err
}

// Delete deletes a message from
func (rabbitMQQueue *RabbitMQQueue) Delete(ctx context.Context, receiptHandle string) error {
	// Get the delivery tag from the receipt handle
	deliveryTag, err := strconv.ParseUint(receiptHandle, 10, 64)
	if err != nil {
		return err
	}
	// Reject a single message
	return rabbitMQQueue.channel.Reject(deliveryTag, false)
}
