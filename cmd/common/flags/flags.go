package flags

import "flag"

var Broker = flag.String("broker", "127.0.0.1:9092", "kafka broker to consume messages")
var Topic = flag.String("topic", "currencyUpdate", "kafka topic to consume messages")
var Debug = flag.Bool("debug", false, "debug mode")
var Trace = flag.Bool("trace", false, "trace mode")
