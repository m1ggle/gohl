.PHONY: build test run clean install help

# 项目变量
BINARY_NAME=gohl
MAIN_PATH=./main.go

# 获取当前系统信息
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# 默认帮助信息
help: ## 显示帮助信息
	@echo "使用方法: make <目标>"
	@echo ""
	@echo "可用的目标:"
	@grep -E '^[a-zA-Z_0-9%-]+:.*?## .*$$' $(word 1,$(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## 构建二进制文件
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

test: ## 运行测试
	go test ./...

test-v: ## 运行详细测试
	go test -v ./...

install: ## 安装到 GOPATH/bin
	go install $(MAIN_PATH)

run: ## 运行程序
	go run $(MAIN_PATH)

clean: ## 清理构建产物
	rm -f bin/$(BINARY_NAME)

# 跨平台构建
build-linux: ## 构建 Linux 版本
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-mac: ## 构建 macOS 版本
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

build-windows: ## 构建 Windows 版本
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

build-all: build-linux build-mac build-windows ## 构建所有平台版本

fmt: ## 格式化代码
	go fmt ./...

vet: ## 检查代码错误
	go vet ./...

lint: ## 代码静态检查
	golangci-lint run ./... || echo "请安装 golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

tidy: ## 更新依赖
	go mod tidy

check: fmt vet test ## 执行所有检查
