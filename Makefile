vendor:
	@glide install -v

install:
	@go install ./cmd/ci

run: install
	@ci

run-k: install
	@KUBERNETES="k8s" ci --cluster kubernetes

clean:
	@./clean.sh

.PHONY: run run-k build vendor clean
