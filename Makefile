BIN := "./bin/fileprocess"
FILE := "./reverse-duplicates.txt"
NPROC := 8

build:
	go build -o $(BIN) ./main.go

run: build
	$(BIN) $(NPROC) $(FILE)

.PHONY: build run
