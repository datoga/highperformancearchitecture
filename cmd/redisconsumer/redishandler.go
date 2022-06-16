package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/datoga/highperformancearchitecture"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type RedisHandler struct {
	client redis.UniversalClient
	debug  bool
}

func CreateRedisClient() redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

func NewRedisHandler(client redis.UniversalClient) *RedisHandler {
	return &RedisHandler{client: client}
}

func (rh *RedisHandler) WithDebug() *RedisHandler {
	rh.debug = true

	return rh
}

func (rh *RedisHandler) Handle(msg *message.Message) error {
	var currencyUpdate highperformancearchitecture.CurrencyUpdate

	if err := json.Unmarshal(msg.Payload, &currencyUpdate); err != nil {
		fmt.Println("Failed unmarshalling message:", err)
		return err
	}

	key := currencyUpdate.Exchange.Origin + "_" + currencyUpdate.Exchange.Target
	value := currencyUpdate.Exchange.Rate

	err := rh.client.Set(
		context.Background(),
		key,
		value,
		0,
	).Err()

	if err != nil {
		fmt.Println("Failed setting key:", err)
		if errors.Is(err, redis.ErrClosed) {
			_ = rh.client.Close()
			rh.client = CreateRedisClient()
		}

		return err
	}

	if rh.debug {
		fmt.Printf("Key %s updated with value %f\n", key, value)
	}

	return nil
}
