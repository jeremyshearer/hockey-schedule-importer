BIN = hockey-schedule-importer

.PHONY: build test clean run serve

build:
	go build -o $(BIN) .

test:
	go test ./...

clean:
	rm -f $(BIN)

run: build
	./$(BIN) convert

serve: build
	./$(BIN) httpd
