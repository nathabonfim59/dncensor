.PHONY: all build build-dev vet tidy clean

all: build vet

build:
	go build -o dncensor .

build-dev:
	go build -tags dev -o dncensor-dev .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -f dncensor dncensor-dev
