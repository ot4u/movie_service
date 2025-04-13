all: build

build:
	go mod download
	go build -o ./bin/movie_service

rebuild:
	make build

run: build
	./bin/movie_service

test:
	go test -v ./pool
	make clean

style:
	go fmt ./...

clean:
	rm -rf ./bin *.txt pool/*.txt movie_service