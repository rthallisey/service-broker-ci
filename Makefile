vendor:
	@glide install -v

install:
	@go build -i ./cmd/ci

run: install
	@ci

run-k: install
	@ci --cluster kubernetes

.PHONY: run run-k build vendor
