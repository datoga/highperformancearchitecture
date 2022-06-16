package main

import (
	"github.com/datoga/highperformancearchitecture/cmd/common/flags"
	"github.com/datoga/highperformancearchitecture/cmd/common/stats"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/biter777/countries"
)

var currencies = countries.AllCurrencies()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	logger := watermill.NewStdLogger(*flags.Debug, *flags.Trace)

	statsManager := stats.New()

	go statsManager.Run()

	sleep, err := time.ParseDuration(*pace)

	if err != nil {
		panic(err)
	}

	pub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{*flags.Broker},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	producer := NewPublisher(pub, *flags.Topic)

	if *flags.Debug {
		producer = producer.WithDebug()
	}

	defer pub.Close()

	for {
		go func() {
			if errPublish := producer.PublishOne(); errPublish != nil {
				panic(errPublish)
			}
		}()

		statsManager.CountMessage()

		time.Sleep(sleep)
	}
}
