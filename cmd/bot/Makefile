deps:
	go get
run:
	go run main.go api.go helpers.go
build:
	go build -o ../../bin/bot main.go api.go helpers.go
build-docker:
	env GOOS=linux GOARCH=amd64 go build -o ../../bin/bot main.go api.go helpers.go
