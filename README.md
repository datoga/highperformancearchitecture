[![Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/datoga/highperformancearchitecture)

# High performance event driven architecture using kafka and go

Sample project with a high performance architecture

## Tools and libs used

- **Go**: https://golang.org.
- **Redpanda**: high-performance kafka compatible replacement https://redpanda.com.
- **Kafkacat**: utility to produce and consume messages from kafka.
- **Redis**: in-memory data store https://redis.io.
- **Watermill**: library to work with events in Go: https://watermill.io.
- **JSON iterator**: high-performance 100% compatible drop-in replacement of "encoding/json" https://github.com/json-iterator/go.
- **Gitpod**: to launch the demo: https://gitpod.io.

## Step 1: Producing events and consuming with kafkacat

```
./producer -trace
kafkacat -b localhost:9092 -t updateCurrency
```

Check with the param `-pace <duration>` the producing rate, por example:

```
./producer -pace 1ms
./producer -pace 1us
```

## Step 2: Consuming events with simple consumer (kafkacat replacement)

```
./producer -debug
./simpleconsumer -debug
```

```
./producer -pace 1us
./simpleconsumer
```

## Step 3: Consuming events with redis consumer

```
./producer -debug
./redisconsumer -debug
redis-cli monitor
```

```
./producer -pace 1ms
./redisconsumer
redis-cli monitor
```

One consumer is not enough under high load conditions:
```
./producer -pace 1us
./redisconsumer
```

### Improvements

1. Buffering with redis (and pipeline usage).
2. Configure multiple partitions and multiple consumers.