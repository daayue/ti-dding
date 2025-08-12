# 钉钉群管理工具 Makefile

# 变量定义
BINARY_NAME=ti-dding
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# 默认目标
.PHONY: all
all: build

# 构建应用
.PHONY: build
build:
	@echo "构建 ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/main.go
	@echo "构建完成: ${BUILD_DIR}/${BINARY_NAME}"

# 安装到系统
.PHONY: install
install:
	@echo "安装 ${BINARY_NAME}..."
	go install ${LDFLAGS} ./cmd/main.go
	@echo "安装完成"

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf ${BUILD_DIR}
	@go clean
	@echo "清理完成"

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	go test -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 检查代码质量
.PHONY: lint
lint:
	@echo "检查代码质量..."
	golangci-lint run

# 更新依赖
.PHONY: deps
deps:
	@echo "更新依赖..."
	go mod tidy
	go mod download

# 运行应用
.PHONY: run
run: build
	@echo "运行 ${BINARY_NAME}..."
	@${BUILD_DIR}/${BINARY_NAME} --help

# 创建示例配置
.PHONY: config
config:
	@echo "创建示例配置文件..."
	@cp configs/config.yaml.example configs/config.yaml
	@echo "请编辑 configs/config.yaml 文件，填入你的钉钉应用信息"

# 创建示例数据
.PHONY: data
data:
	@echo "创建示例数据文件..."
	@mkdir -p data
	@echo "示例CSV文件已创建: data/groups_example.csv"

# 帮助信息
.PHONY: help
help:
	@echo "钉钉群管理工具 Makefile"
	@echo ""
	@echo "可用目标:"
	@echo "  build          - 构建应用"
	@echo "  install        - 安装到系统"
	@echo "  clean          - 清理构建文件"
	@echo "  test           - 运行测试"
	@echo "  test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  fmt            - 格式化代码"
	@echo "  lint           - 检查代码质量"
	@echo "  deps           - 更新依赖"
	@echo "  run            - 运行应用"
	@echo "  config         - 创建示例配置文件"
	@echo "  data           - 创建示例数据文件"
	@echo "  help           - 显示此帮助信息"
	@echo ""
	@echo "示例用法:"
	@echo "  make build     # 构建应用"
	@echo "  make install   # 安装到系统"
	@echo "  make config    # 创建配置文件"
	@echo "  make run       # 运行应用" 