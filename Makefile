.PHONY: all build build-dev vet tidy clean

all: build vet

build:
	go build -o build/dncensor .

build-dev:
	go build -tags dev -o build/dncensor-dev .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf build/
