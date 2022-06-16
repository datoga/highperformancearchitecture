package stats

import (
	"fmt"
	"time"
)

const defaultPeriod = 5

type Stats struct {
	periodInSeconds uint
	totalMessages   uint
	messages        uint
}

func New() Stats {
	return Stats{periodInSeconds: defaultPeriod}
}

func (s *Stats) CountMessage() {
	s.totalMessages++
	s.messages++
}

func (s *Stats) Run() {
	for {
		time.Sleep(time.Duration(s.periodInSeconds) * time.Second)

		fmt.Printf("\n\n")
		fmt.Printf("----------------------------------------------\n")
		fmt.Printf("Total messages: %v, Rate: %.02f messages/s\n", s.totalMessages, float64(s.messages)/float64(s.periodInSeconds))
		fmt.Printf("----------------------------------------------\n")
		fmt.Printf("\n\n")

		s.messages = 0
	}
}
