package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/datoga/highperformancearchitecture/cmd/common/flags"
	"github.com/datoga/highperformancearchitecture/cmd/common/stats"
)

var statsManager stats.Stats

func main() {
	statsManager = stats.New()

	go statsManager.Run()

	logger := watermill.NewStdLogger(*flags.Debug, *flags.Trace)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		middleware.InstantAck,
		middleware.Recoverer,
	)

	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       []string{*flags.Broker},
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: "redis",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	internalTopic := "internal"

	_ = router.AddHandler(
		"route_to_internal",
		*flags.Topic,
		sub,
		internalTopic,
		pubSub,
		func(msg *message.Message) ([]*message.Message, error) {
			return message.Messages{msg}, nil
		},
	)

	if *flags.Debug {
		router.AddNoPublisherHandler(
			"print_incoming_messages",
			internalTopic,
			pubSub,
			printMessages,
		)
	}

	rdb := CreateRedisClient()

	defer rdb.Close()

	redisHandler := NewRedisHandler(rdb)

	if *flags.Debug {
		redisHandler = redisHandler.WithDebug()
	}

	router.AddNoPublisherHandler(
		"update_redis",
		internalTopic,
		pubSub,
		redisHandler.Handle,
	)

	router.AddNoPublisherHandler(
		"count_messages",
		internalTopic,
		pubSub,
		countMessage,
	)

	if err = router.Run(context.Background()); err != nil {
		panic(err)
	}
}

func countMessage(_ *message.Message) error {
	statsManager.CountMessage()

	return nil
}

func printMessages(msg *message.Message) error {
	fmt.Printf(
		"\n> Received message: %s\n> %s\n> metadata: %v\n\n",
		msg.UUID, string(msg.Payload), msg.Metadata,
	)
	return nil
}
