all: producer simpleconsumer redisconsumer

producer:
	go build -o bin/producer ./cmd/producer

simpleconsumer:
	go build -o bin/simpleconsumer ./cmd/simpleconsumer

redisconsumer:
	go build -o bin/redisconsumer ./cmd/redisconsumer