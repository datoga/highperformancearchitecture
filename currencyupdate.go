package highperformancearchitecture

type CurrencyUpdate struct {
	ID        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	Exchange  Exchange `json:"exchange"`
}

type Exchange struct {
	Origin string  `json:"origin"`
	Target string  `json:"target"`
	Rate   float64 `json:"rate"`
}
