deps:
	go get
run:
	go run main.go
build:
	go build -o ../../bin/server main.go
build-docker:
	env GOOS=linux GOARCH=amd64 go build -o ../../bin/server main.go