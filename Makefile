vendor:
	@glide install -v
run:
	@go build -i -ldflags="-s -w" ./cmd/ci
	@./ci

.PHONY: run vendor
