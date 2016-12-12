.PHONY: pb build test deps

pb:
	for f in pb/**/*.proto; do \
		protoc --go_out=plugins=grpc:. $$f; \
		echo compiled: $$f; \
	done

deps:
	go get ./...

test: 
	go test ./...

build: deps test
	@./bin/build.sh

clean:
	@./bin/clean.sh

up:
	docker-compose build
	docker-compose up -d

down:
	docker-compose down 
