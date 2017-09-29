vendor:
	@glide install -v

build:
	@go build -i -ldflags="-s -w" ./cmd/ci

run: build
	@./ci

run-k: build
	@./ci --cluster kubernetes

.PHONY: run run-k build vendor
