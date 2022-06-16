package main

import (
	"context"
	"fmt"
	"github.com/datoga/highperformancearchitecture/cmd/common/flags"
	"github.com/datoga/highperformancearchitecture/cmd/common/stats"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

var (
	messages = 0
	initTime = time.Now()
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	logger := watermill.NewStdLogger(*flags.Debug, *flags.Trace)

	statsManager := stats.New()

	go statsManager.Run()

	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       []string{*flags.Broker},
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: "simple",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	defer sub.Close()

	chMessage, err := sub.Subscribe(context.Background(), *flags.Topic)

	if err != nil {
		panic(err)
	}

	for message := range chMessage {
		if *flags.Debug {
			fmt.Println(string(message.Payload))
		}
		message.Ack()

		statsManager.CountMessage()
	}
}
