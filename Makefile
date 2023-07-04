.PHONY: build
build:
	docker build -f docker/client/Dockerfile -t ddos-pow-client:latest .
	docker build -f docker/server/Dockerfile -t ddos-pow-server:latest .

.PHONY: run
run:
	@echo 'Running locally...'
	docker-compose up --build

.PHONY: test
test:
	@echo 'Running test...'
	go test -cover -count=1 -v ./...
