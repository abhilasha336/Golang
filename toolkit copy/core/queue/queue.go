package queue

import "context"

// Queue defines methods for working with various queue providers.
type Queue interface {
	// Build constructs a message in a format suitable for sending to the queue.
	// The specific implementation of the message construction may vary based on the
	// type of queue provider used. The returned interface{} is the
	// constructed message, and error is returned if there was any issue in message construction.
	ComposeMessage(ctx context.Context, message []byte) (interface{}, error)

	// Send sends a message to a queue. Depending on the type of queue used, the input
	// parameters may vary. For example, if RabbitMQ is used, the input must be of the form
	// amqp.Publishing, which is defined in the package "github.com/rabbitmq/amqp091-go".
	Send(ctx context.Context, message interface{}) error

	// Receive retrieves data from the queue. Similar to the Send method, the output may vary based
	// on the type of queue used. If RabbitMQ is used, the output is of the form <-chan amqp.Delivery,
	// which is defined in the package "github.com/rabbitmq/amqp091-go".
	Receive(ctx context.Context) (interface{}, error)

	// Delete deletes a message from the queue. The implementation of this method may vary depending
	// on the type of queue used. For RabbitMQ, To delete a message
	// you need to provide the delivery tag of the message to be deleted.
	Delete(ctx context.Context, receiptHandle string) error

	// Close closes a connection to the queue.
	Close() error
}
