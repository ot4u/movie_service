build:
	docker-compose build

run:
	docker-compose up --build

stop:
	docker-compose down

restart:
	docker-compose restart

clean:
	docker-compose down -v

fmt:
	go fmt ./...

install:
	go mod tidy