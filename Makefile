up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: up

down:
	docker-compose down --remove-orphans
.PHONY: down