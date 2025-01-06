package config

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection
var RabbitMQCh *amqp.Channel

// SetupRabbitMQ initializes the global RabbitMQ connection and channel
func SetupRabbitMQ() {
	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	fmt.Println(addr)

	var err error
	RabbitMQConn, err = amqp.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	RabbitMQCh, err = RabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}

	log.Println("RabbitMQ connection and channel initialized.")
}

// InitQueue sets up a durable queue
func InitQueue(queueName string) {
	_, err := RabbitMQCh.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	log.Printf("Queue declared: %s", queueName)
}

// BindQueueToExchange binds a queue to an exchange with a specific routing key
func BindQueueToExchange(queueName, exchangeName, routingKey string) {
	err := RabbitMQCh.QueueBind(
		queueName,    // Queue name
		routingKey,   // Routing key
		exchangeName, // Exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue %s to exchange %s: %v", queueName, exchangeName, err)
	}
	log.Printf("Queue %s bound to exchange %s with routing key %s", queueName, exchangeName, routingKey)
}

// InitDirectRabbitMQExchange sets up a direct exchange for notifications
func InitDirectRabbitMQExchange(exchangeName string) {
	err := RabbitMQCh.ExchangeDeclare(
		exchangeName, // Exchange name
		"direct",     // Type
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare direct exchange %s: %v", exchangeName, err)
	}
	log.Printf("Declared RabbitMQ direct exchange: %s", exchangeName)
}

func InitFanoutRabbitMQExchange(exchangeName string) {
	err := RabbitMQCh.ExchangeDeclare(
		exchangeName, // Exchange name
		"fanout",     // Type
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare fanout exchange %s: %v", exchangeName, err)
	}
	log.Printf("Declared RabbitMQ fanout exchange: %s", exchangeName)
}

// CleanupRabbitMQ closes the RabbitMQ connection and channel
func CleanupRabbitMQ() {
	if RabbitMQCh != nil {
		if err := RabbitMQCh.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}
	if RabbitMQConn != nil {
		if err := RabbitMQConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}
	log.Println("RabbitMQ connection and channel closed.")
}

// queueExists checks if a RabbitMQ queue exists
func QueueExists(queueName string) bool {
	_, err := RabbitMQCh.QueueDeclarePassive(
		queueName, // Queue name
		true,      // Durable
		true,      // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Printf("Queue %s does not exist: %v", queueName, err)
		return false
	}
	return true
}

func CleanupQueue(queueName string) error {
	_, err := RabbitMQCh.QueueDelete(
		queueName, // Queue name
		false,     // IfUnused
		false,     // IfEmpty
		false,     // NoWait
	)
	return err
}