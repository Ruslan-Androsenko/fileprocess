BIN := "./bin/fileprocess"

build:
	go build -o $(BIN) ./main.go

run: build
	$(BIN)

.PHONY: build run
