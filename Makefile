test: ## To test .rego files, run make test.
	@opa test . -v

up:
	@docker-compose up -d

down:
	@docker-compose down

restart: down up

server:
	@go run server.go

local:
	@opa run --server --bundle example
