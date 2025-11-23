.PHONY: all frontend backend clean run

# Default target
all: frontend backend

# Build the Vue.js frontend
frontend:
	@echo "Building Frontend..."
	cd wsjtx-web-ui && npm install && npm run build

# Build the Go backend
backend:
	@echo "Building Backend..."
	go build -o wsjtx-web main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f wsjtx-web
	rm -rf wsjtx-web-ui/dist

# Run the application
run: all
	@echo "Starting application..."
	./wsjtx-web
