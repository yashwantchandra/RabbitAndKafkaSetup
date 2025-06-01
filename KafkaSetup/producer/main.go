package main

import (
    "context"
    "fmt"
    "log"

    "github.com/segmentio/kafka-go"
)

func main() {
    // Kafka writer configuration
    writer := &kafka.Writer{
        Addr:     kafka.TCP("localhost:9092"),
        Topic:    "test-topic-with-2-partition",
       // Balancer: &kafka.Hash{},
    }

    defer func() {
        if err := writer.Close(); err != nil {
            log.Fatalf("Failed to close Kafka writer: %v", err)
        }
    }()

    fmt.Println("Starting Kafka producer...")

    // Produce messages to both partitions
    for i := 0; i < 10000; i++ {
        partition := i % 2 // Alternate between partition 0 and 1
        message := kafka.Message{
            Key:       []byte(fmt.Sprintf("key-%d", i)),
            Value:     []byte(fmt.Sprintf("message-%d", i)),
            Partition: partition, // Explicitly set the partition
        }

        err := writer.WriteMessages(context.Background(), message)
        if err != nil {
            log.Printf("Failed to write message to partition %d: %v", partition, err)
        } else {
            fmt.Printf("Produced message to partition %d: %s\n", partition, message.Value)
        }
    }
}