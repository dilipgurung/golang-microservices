.PHONY: pb build test deps

pb:
	n pb/**/*.proto; do \
		protoc --go_out=plugins=grpc:. $$f; \
		echo compiled: $$f; \
	done

deps:
	go get

build:
	./bin/build.sh

up:
	docker-compose build
	docker-compose up -d

down:
	docker-compose down 
