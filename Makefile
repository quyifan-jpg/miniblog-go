.PHONY: help build run test clean docker-build docker-up docker-down install

help: ## 显示帮助信息
	@echo "Simple IM Makefile Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装依赖
	go mod download
	go mod tidy

build: ## 编译项目
	go build -o bin/miniblog cmd/server/main.go cmd/server/http.go cmd/server/websocket.go

run: ## 运行项目
	go run cmd/server/main.go cmd/server/http.go cmd/server/websocket.go

test: ## 运行测试
	go test -v ./...

clean: ## 清理构建文件
	rm -rf bin/
	rm -rf logs/

docker-build: ## 构建 Docker 镜像
	docker build -t miniblog:latest -f docker/Dockerfile .

docker-up: ## 启动 Docker Compose
	docker-compose -f docker/docker-compose.yml up -d

docker-down: ## 停止 Docker Compose
	docker-compose -f docker/docker-compose.yml down

docker-logs: ## 查看 Docker 日志
	docker-compose -f docker/docker-compose.yml logs -f

lint: ## 运行代码检查
	golangci-lint run

fmt: ## 格式化代码
	go fmt ./...

dev: ## 开发模式运行
	air

