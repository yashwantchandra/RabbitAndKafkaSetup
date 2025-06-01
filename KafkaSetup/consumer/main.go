package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
)

func main() {
    fmt.Println("Starting Kafka consumer...")

    // Start one goroutine for each partition
    go consumePartition(0)
    go consumePartition(1)

    // Block forever
    select {}
}

func consumePartition(partition int) {
    for {
        fmt.Printf("Starting consumer for partition %d...\n", partition)

        reader := kafka.NewReader(kafka.ReaderConfig{
            Brokers:   []string{"localhost:9092"},
            Topic:     "test-topic-with-2-partition",
            Partition: partition,
           // GroupID:   "my-group", // Not necessary when using Partition
            MaxWait:   1 * time.Second,
        })
		reader.SetOffset(kafka.LastOffset) // Start reading new contents only without this consumer will // read all the messages from the beginning 
		
		// if  we provide group id then kafka will remember the offset of the last read message and will start reading from that offset next time we run the consumer

        // Wrap the read loop with a recover block
        func() {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("Recovered from panic in partition %d reader: %v", partition, r)
                }
                reader.Close() // Ensure reader is closed before restarting
            }()

            for {
                msg, err := reader.ReadMessage(context.Background())
                if err != nil {
                    log.Printf("Error reading from partition %d: %v", partition, err)
                    // Optionally add a backoff before retrying
                    time.Sleep(2 * time.Second)
                    return // Exit loop to restart the goroutine
                }

                fmt.Printf("Received from partition %d: %s\n",msg.Partition, msg.Value)
            }
        }()

        // Short delay before restarting the goroutine
        log.Printf("Restarting reader for partition %d in 3 seconds...\n", partition)
        time.Sleep(3 * time.Second)
    }
}