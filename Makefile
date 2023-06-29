.PHONY: build
build:
	docker build -f cmd/client/Dockerfile -t ddos-pow-client:latest .
	docker build -f cmd/server/Dockerfile -t ddos-pow-server:latest .

.PHONY: local-run
local-run:
	@echo 'Running local...'
	docker-compose -p test up

.PHONY: test
test:
	@echo 'Running test...'
	go test -race -cover -count=1 -v ./...
