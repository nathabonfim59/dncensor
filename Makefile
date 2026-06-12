BINARY := dncensor

.PHONY: all build vet tidy clean run current list

all: tidy vet build

build:
	go build -o $(BINARY) .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)

run: build
	sudo ./$(BINARY)

current:
	./$(BINARY) current

list:
	./$(BINARY) list-providers
