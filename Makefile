.PHONY: build-production
build-production:
	docker buildx build --platform linux/amd64 -t go-boilerplate-amd64 .

.PHONY: build
build:
	docker build -t go-boilerplate .

.PHONY: reset
reset:
	docker compose down -v

.PHONY: run
run:
	@echo 'Running in hot-reload mode'
	docker compose --profile dev up --build --watch


.PHONY: migration-down
migration-down:
	docker run -v ./database/migrations:/migrations --network host migrate/migrate -path /migrations/ -database postgres://bot:admin@localhost:5432/master?sslmode=disable down -all
