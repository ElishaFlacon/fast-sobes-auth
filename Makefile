.PHONY: run run-dev run-prod down stop clean clean-all logs update lint tidy deps-reset deps-upgrade deps-cleancache

# ==============================================================================
# Docker compose commands

run:
	@echo "Starting docker containers..."
	docker compose -f docker-compose.yml up --build

run-dev:
	@echo "Starting dev docker containers..."
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

run-prod:
	@echo "Starting prod docker containers..."
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# ==============================================================================
# Docker support

down:
	@echo "Stopping and removing all docker containers"
	docker compose down

stop:
	@echo "Stopping docker containers"
	docker compose stop

clean:
	@echo "Cleaning docker data..."
	docker system prune -f

# DO NOT USE IF YOU DONT KNOW WHAT IS IT
clean-all:
	@echo "Cleaning ALL docker data..."
	docker system prune -a --volumes -f

logs:
	@echo "View docker containers logs..."
	docker compose logs -f

# ==============================================================================
# Tools commands

lint:
	@echo "Starting linters"
	golangci-lint run ./...

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache
