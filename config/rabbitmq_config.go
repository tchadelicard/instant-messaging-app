package config

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)


func SetupRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	fmt.Println(addr)

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}

	log.Println("RabbitMQ connection and channel initialized.")
	return conn, ch
}

// InitFanOutRabbitMQBindings sets up a fan-out exchange for a queue
func InitFanOutRabbitMQBindings(ch *amqp.Channel, queueName, exchangeName string) {
	// Declare fan-out exchange
	err := ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,  // Durable
		false, // Auto-deleted
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare a unique queue
	q, err := ch.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Delete when unused
		true,  // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,     // Queue name
		"",         // Routing key (not used in fan-out)
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	log.Printf("Fan-out bindings initialized for queue: %s", q.Name)
}

// InitDirectRabbitMQExchange sets up a direct exchange for notifications
func InitDirectRabbitMQExchange(ch *amqp.Channel, exchangeName string) {
	err := ch.ExchangeDeclare(
		exchangeName, // Exchange name
		"direct",     // Type
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}
	log.Printf("Declared RabbitMQ direct exchange: %s", exchangeName)
}