.PHONY: build-production
build-production:
	docker buildx build --platform linux/amd64 -t go-boilerplate-amd64 .

.PHONY: build
build:
	docker build -t go-boilerplate .

.PHONY: reset
reset:
	docker compose --profile dev down -v --remove-orphans

.PHONY: run
run:
	@echo 'Running in hot-reload mode'
	docker compose --profile dev up --build --watch

.PHONY: smoke
smoke:
	go test -tags=smoke -count=1 -v ./test/smoke/...

.PHONY: migration-down
migration-down:
	docker run -v ./database/migrations:/migrations --network host migrate/migrate -path /migrations/ -database postgres://bot:admin@localhost:5432/master?sslmode=disable down -all

.PHONY: deploy
deploy:
	ansible-playbook -i etc/ansible/inventory.ini etc/ansible/playbooks/deploy.yml
