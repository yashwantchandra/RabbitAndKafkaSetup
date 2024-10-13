package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}
func main() {
    // Step 1: Connect to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    // Step 2: Create a channel
    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    // Step 3: Declare a queue (same as in Producer)
    // q, err := ch.QueueDeclare(
    //     "test_queue", // queue name
    //     false,        // durable
    //     false,        // delete when unused
    //     false,        // exclusive
    //     false,        // no-wait
    //     nil,          // arguments
    // )
    // failOnError(err, "Failed to declare a queue")

    // Step 4: Consume messages from the queue

    // Step 5: Create a Go channel to receive the message
    forever := make(chan bool)

    go func() {
        for {
            msgs, err := ch.Consume(
                "query", // queue
                "",     // consumer
                false,   // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
            )
            if err != nil {
                log.Printf("Error consuming messages, retrying: %s", err)
                time.Sleep(5 * time.Second) // Wait before retrying
                continue
            }
    
            // Consume the messages
            for d := range msgs { // this loop will terminate if the channel is closed 
                log.Printf("Received a message: %s", d.Body)
                d.Ack(false) // Acknowledge message
            }
    
            log.Println("Channel closed, attempting to reconnect...")
        }
    }()
    
    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}