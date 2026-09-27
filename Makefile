BINARY_DIR := ./bin
CMD_DIR := ./cmd/...

.PHONY: all build clean run-server run-transcoder

all: clean build

build:
	@echo "🔨 Building all binaries..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/ $(CMD_DIR)
	@echo "✅ Build complete! $(BINARY_DIR)/"

clean:
	@echo "🧹 Cleaning up binaries..."
	@rm -rf $(BINARY_DIR)
	@echo "✨ Clean complete!"

run-server: build
	@echo "🚀 Running server..."
	@$(BINARY_DIR)/server

## run-cli: Quick shortcut to run a specific binary (e.g., cli)
run-transcoder: build
	@echo "🚀 Running transcoder..."
	@$(BINARY_DIR)/transcoder
