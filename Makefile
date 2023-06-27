.PHONY: build
build:
	docker build -t ddos-pow:latest .

.PHONY: local-run
local-run:
	@echo 'Running local...'
	docker-compose -p test up
