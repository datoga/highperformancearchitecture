package main

import (
	"encoding/json"
	"fmt"
	"github.com/datoga/highperformancearchitecture"
	"github.com/datoga/highperformancearchitecture/cmd/common/flags"
	"github.com/datoga/highperformancearchitecture/cmd/common/stats"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/biter777/countries"
	"github.com/oklog/ulid/v2"
)

var (
	entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
)

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

	defer pub.Close()

	currencies := countries.AllCurrencies()

	for {
		go sendMessage(currencies, pub)

		statsManager.CountMessage()

		time.Sleep(sleep)
	}
}

func sendMessage(currencies []countries.CurrencyCode, pub *kafka.Publisher) {
	currency1 := pickCurrency(currencies)
	currency2 := pickCurrency(currencies)

	rate := rand.Float64() * 10

	msgULID := watermill.NewULID()

	currencyUpdate := highperformancearchitecture.CurrencyUpdate{
		ID:        msgULID,
		Timestamp: time.Now().UnixNano(),
		Exchange: highperformancearchitecture.Exchange{
			Origin: currency1.Alpha(),
			Target: currency2.Alpha(),
			Rate:   rate,
		},
	}

	bytes, errMarshal := json.Marshal(currencyUpdate)

	if errMarshal != nil {
		panic(errMarshal)
	}

	msg := message.NewMessage(msgULID, bytes)

	if *flags.Debug {
		fmt.Printf("Publishing %+v\n", currencyUpdate)
	}

	if errPublish := pub.Publish("currencyUpdate", msg); errPublish != nil {
		panic(errPublish)
	}
}

func pickCurrency(currencies []countries.CurrencyCode) countries.CurrencyCode {
	randomIndex := rand.Intn(len(currencies))
	return currencies[randomIndex]
}
