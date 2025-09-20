.PHONY: build run clean test deps swagger

# 构建项目
build:
	go build -o bin/glimgate main.go

# 运行项目
run:
	go run main.go

# 清理构建文件
clean:
	rm -rf bin/

# 运行测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod tidy
	go mod download

# 生成Swagger文档
swagger:
	swag init

# 开发环境运行（带热重载）
dev:
	air

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 初始化项目
init: deps swagger

build-compose:
	docker compose build

up:
	docker compose up -d

deploy:  swagger build-compose up
	@echo "部署完成"

down:
	docker compose down

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  build    - 构建项目"
	@echo "  run      - 运行项目"
	@echo "  clean    - 清理构建文件"
	@echo "  test     - 运行测试"
	@echo "  deps     - 安装依赖"
	@echo "  swagger  - 生成Swagger文档"
	@echo "  dev      - 开发环境运行（需要安装air）"
	@echo "  fmt      - 格式化代码"
	@echo "  lint     - 代码检查（需要安装golangci-lint）"
	@echo "  init     - 初始化项目"
	@echo "  help     - 显示帮助信息"