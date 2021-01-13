deps:
	go get
run:
	go run main.go api.go helpers.go
build:
	go build -o ../../bin/bot main.go api.go helpers.go
build-docker:
	env GOOS=linux GOARCH=amd64 go build -o ../../bin/bot main.go api.go helpers.go

# env ADDRESS=127.0.0.1 KEYSPACE=casspoll make run