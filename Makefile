APP_NAME = biathlon

build:
	go build -o $(APP_NAME) ./cmd/main.go

run: build
	./$(APP_NAME)

clean:
	rm -f $(APP_NAME)

all: build