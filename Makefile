PROJECT=github.com/crazylazyowl/metrics-tpl

build: #### Build the server and agent binaries.
	go build -o ./bins/server ./cmd/server
	go build -o ./bins/agent ./cmd/agent

tools: #### Install all necessary tools.
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	go get github.com/golang/mock



test_mock: #### Generate mock interfaces.
	mockgen -destination=./internal/controller/httprest/api/mocks/ping_mock.go -package=mocks $(PROJECT)/internal/usecase/ping Pinger


test: #### Run unit tests.
	go test -race -v ./...


lint: #### Run linter.
	golangci-lint run


help: #### Show this help message.
	@sed -e '/__hidethis__/d; /###/!d; s/:.\+#### /\t\t/g; s/:.\+#### /\t\t\t/g; s/:.\+### /\t/g' $(MAKEFILE_LIST)
