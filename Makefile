all: lint run

.PHONY: link
lint:
	golint -min_confidence=0.3 ./...
	staticcheck ./...

.PHONY: run
run:
	go build -o twinui .
	./twinui -port 3000 -story ./model/gopher.json
