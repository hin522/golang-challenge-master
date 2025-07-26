
DATABASE = docker compose run --rm database
GOLANG = docker compose run --rm --service-ports golang
GOLANG_IMAGE = golang:1.24.5-alpine

RED = \033[1;31m
CYAN = \033[0;36m

.PHONY: startdb
startdb:
	docker compose up database
	echo "$(CYAN)Database exited"

.PHONY: rebuilddb
rebuilddb:
	docker compose down -v database 
	docker compose build database
	echo "$(CYAN)Database rebuilt successfully."

.PHONY: run
run: 
	docker compose up --build
	echo "$(CYAN)Server started"

.PHONY: build
build: 
	docker run --rm -v ${PWD}:/app -w /app $(GOLANG_IMAGE) \
	    sh -c "apk add --no-cache git && go build -o bin/server ./server/"
	echo "$(CYAN)Server built successfully."
