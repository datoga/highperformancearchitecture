package stats

import (
	"fmt"
	"time"
)

type Stats struct {
	messages uint
	init     time.Time
}

func New() Stats {
	return Stats{
		init: time.Now(),
	}
}

func (s *Stats) CountMessage() {
	s.messages++
}

func (s *Stats) Run() {
	for {
		time.Sleep(time.Second)

		now := time.Now()

		elapsed := now.Sub(s.init)

		fmt.Printf("\n\n")
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("Elapsed: %v, Messages: %v, Rate: %.02f messages/s\n", elapsed, s.messages, float64(s.messages)/elapsed.Seconds())
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("\n\n")
	}
}
