APP_NAME = giproxy
ENTRY_FILE = ./cmd/main.go

install:
	@echo ">>> Installing dependencies..."
	go mod tidy
	go mod download


install-nodemon:
	@echo ">>> Checking if nodemon is installed..."
	@if ! [ -x "$$(command -v nodemon)" ]; then \
		echo ">>> Installing nodemon globally with npm..."; \
		npm install -g nodemon; \
	else \
		echo ">>> nodemon is already installed."; \
	fi

run:
	@echo ">>> Starting application..."
	go run $(ENTRY_FILE)

clean:
	@echo ">>> Cleaning up temporary files..."
	rm -rf ./tmp

watch:
	@echo ">>> Starting application with nodemon..."
	nodemon --exec "go run $(ENTRY_FILE)" --ext go --signal SIGKILL

build:
	@echo ">>> Building application..."
	go build -o ./bin/$(APP_NAME) $(ENTRY_FILE)

reset:
	@echo ">>> Resetting dependencies..."
	rm -rf go.sum
	go mod tidy
	go mod download