package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/datoga/highperformancearchitecture"
	"github.com/datoga/highperformancearchitecture/cmd/common/flags"
	"github.com/go-redis/redis/v8"
)

func main() {
	logger := watermill.NewStdLogger(*flags.Debug, *flags.Trace)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer rdb.Close()

	router.AddNoPublisherHandler(
		"update_redis",
		internalTopic,
		pubSub,
		func(msg *message.Message) error {
			var currencyUpdate highperformancearchitecture.CurrencyUpdate

			if err := json.Unmarshal(msg.Payload, &currencyUpdate); err != nil {
				return err
			}

			key := currencyUpdate.Exchange.Origin + "_" + currencyUpdate.Exchange.Target
			value := currencyUpdate.Exchange.Rate

			//TODO: ASYNC
			err := rdb.Set(
				msg.Context(),
				key,
				value,
				0,
			).Err()

			if err != nil {
				return err
			}

			if *flags.Debug {
				fmt.Printf("Key %s updated with value %f\n", key, value)
			}

			return nil
		},
	)

	if err = router.Run(context.Background()); err != nil {
		panic(err)
	}
}

func printMessages(msg *message.Message) error {
	fmt.Printf(
		"\n> Received message: %s\n> %s\n> metadata: %v\n\n",
		msg.UUID, string(msg.Payload), msg.Metadata,
	)
	return nil
}
