up:
	docker-compose up --build -d
.PHONY: up

down:
	docker-compose down --remove-orphans
.PHONY: down