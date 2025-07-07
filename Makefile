PROJECT=github.com/crazylazyowl/metrics-tpl
MOCKS=./internal/usecase/mocks


build: #### Build the server and agent binaries.
	go build -o ./bins/server ./cmd/server
	go build -o ./bins/agent ./cmd/agent


tools: #### Install all necessary tools.
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	go get github.com/golang/mock
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


test_mockgen: #### Generate mock interfaces.
	mockgen -package=mocks -destination=$(MOCKS)/ping_mock.go \
		$(PROJECT)/internal/usecase/ping Pinger
	mockgen -package=mocks -destination=$(MOCKS)/metrics_mock.go \
		$(PROJECT)/internal/usecase/metrics MetricRegistry,MetricFetcher,MetricUpdater


migration_create:
	migrate create -ext sql -dir ./migrations -seq create_metrics_table


test: #### Run unit tests.
	go test -race -v ./...


lint: #### Run linter.
	golangci-lint run --disable=errcheck


help: #### Show this help message.
	@sed -e '/__hidethis__/d; /###/!d; s/:.\+#### /\t\t/g; s/:.\+#### /\t\t\t/g; s/:.\+### /\t/g' $(MAKEFILE_LIST)
