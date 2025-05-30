.PHONY: tidy
tidy:	
	go mod tidy

.PHONY: style
style:
	goimports -l -w ./

.PHONY: unit-test
unit-test:
	go clean -testcache && go test -v ./...

.PHONY: integration-test
integration-test:
	./tests/integration/integration_tests.sh

.PHONY: go-build
go-build:
	CGO_ENABLED=0 go build -o ./bin/flags ./

.PHONY: go-install
go-install:
	go install

.PHONY: build
build:
	docker buildx build --platform linux/amd64 -t github.com/w-h-a/flags:0.1.1-alpha .