vendor:
	@glide install -v

install:
	@go install ./cmd/ci

run: install
	@ci

run-k: install
	@ci --cluster kubernetes

.PHONY: run run-k build vendor
