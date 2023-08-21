package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"os"
	"sync"
)

func main() {
	seeds := []string{"localhost:9093"}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		//kgo.ConsumerGroup("my-group-identifier"),
		kgo.ConsumeTopics("important"),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()
	// 1.) Producing a message
	// All record production goes through Produce, and the callback can be used
	// to allow for synchronous or asynchronous production.
	signal := make(chan os.Signal, 1)
	count := 200
	for {
		//time.Sleep(time.Second)
		select {
		case <-signal:
			fmt.Printf("Program exited")
			os.Exit(1)
		default:
			count++
			var wg sync.WaitGroup
			wg.Add(1)
			record := &kgo.Record{Topic: "important", Value: []byte{byte(count)}}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					fmt.Printf("record had a produce error: %v\n", err)
				} else {
					fmt.Printf("Value %d produced\n", count)
				}
			})
			wg.Wait()

			//// Alternatively, ProduceSync exists to synchronously produce a batch of records.
			//if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
			//	fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
			//} else {
			//	fmt.Printf("Value %d produced\n", count)
			//}
		}
	}
}
