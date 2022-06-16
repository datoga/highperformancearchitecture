package main

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/biter777/countries"
	"github.com/datoga/highperformancearchitecture"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"time"
)

var json = jsoniter.ConfigFastest

//Publisher is an internal publisher of kafka random messages.
type Publisher struct {
	pub   *kafka.Publisher
	topic string
	debug bool
}

//NewPublisher creates an internal publisher for a topic.
func NewPublisher(pub *kafka.Publisher, topic string) *Publisher {
	return &Publisher{
		pub:   pub,
		topic: topic,
	}
}

//WithDebug enables the debug information on publishing data.
func (p *Publisher) WithDebug() *Publisher {
	p.debug = true

	return p
}

//PublishOne publishes one message. It will return an error if the publishing fails.
func (p *Publisher) PublishOne() error {
	currency1 := pickCurrency()
	currency2 := pickCurrency()

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

	bytes, err := json.Marshal(currencyUpdate)

	if err != nil {
		return err
	}

	msg := message.NewMessage(msgULID, bytes)

	if p.debug {
		fmt.Printf("Publishing %+v\n", currencyUpdate)
	}

	return p.pub.Publish(p.topic, msg)
}

func pickCurrency() countries.CurrencyCode {
	randomIndex := rand.Intn(len(currencies))
	return currencies[randomIndex]
}
