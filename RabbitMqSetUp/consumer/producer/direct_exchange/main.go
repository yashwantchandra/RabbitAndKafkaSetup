package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
    if err != nil {
       fmt.Println("%s: %s", msg, err)
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

    // // Step 3: Declare a queue
    // q, err := ch.QueueDeclare(
    //     "test_queue", // queue name
    //     false,        // durable
    //     false,        // delete when unused
    //     false,        // exclusive
    //     false,        // no-wait
    //     nil,          // arguments
    // )
    // failOnError(err, "Failed to declare a queue")

    // Step 4: Publish a message to the queue
    body := "Hello RabbitMQ!"
    err = ch.Publish(
        "Buylead",          // exchange
        "purchase.history",      // routing key (queue name)
        false,       // mandatory
        false,       // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    failOnError(err, "Failed to publish a message")
    fmt.Println(" [x] Sent %s", body)
}