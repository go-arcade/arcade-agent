SHELL := /bin/sh
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

.PHONY: help deps-sync all build run release wire wire-install wire-clean buf buf-install buf-check buf-lint buf-breaking buf-clean

# git
VERSION    = $(shell git describe --tags --always)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# proto
PROTO_DIR := api

LDFLAGS := \
 -X 'github.com/go-arcade/arcade-agent/pkg/version.Version=$(VERSION)' \
 -X 'github.com/go-arcade/arcade-agent/pkg/version.GitBranch=$(GIT_BRANCH)' \
 -X 'github.com/go-arcade/arcade-agent/pkg/version.GitCommit=$(GIT_COMMIT)' \
 -X 'github.com/go-arcade/arcade-agent/pkg/version.BuildTime=$(BUILD_TIME)'

.DEFAULT_GOAL := help

deps-sync: ## 同步Go依赖
	go mod tidy
	go mod verify

help: ## 显示帮助信息
	@echo "Arcade CI/CD 平台 Agent Makefile 命令"
	@echo ""
	@echo "使用方法: make [命令]"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "示例:"
	@echo "  make buf-install    # 安装buf工具"
	@echo "  make buf            # 生成proto代码"
	@echo "  make all            # 完整构建"

all: deps-sync build ## 完整构建

build: wire ## 构建主程序
	go build -ldflags "${LDFLAGS}" -o arcade-agent ./cmd/arcade-agent/

run: deps-sync wire ## 运行主程序（开发模式）
	go run ./cmd/arcade-agent/

release: ## 创建发布版本
	goreleaser --skip-validate --skip-publish --snapshot

# proto代码生成
buf-install: ## 安装buf相关插件
	@echo ">> installing buf..."
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@echo ">> buf installed: $$(which buf)"

buf: buf-check buf-clean ## 生成buf代码
	@echo ">> generating buf code from $(PROTO_DIR)"
	@buf generate buf.build/go-arcade/arcade:main
	@echo ">> buf code generation done."

buf-check: ## 检查buf工具是否已安装
	@command -v buf >/dev/null 2>&1 || { \
		echo "错误: buf 未安装，请先运行 make buf-install"; \
		exit 1; \
	}
	@echo ">> buf installed: $$(which buf)"

buf-lint: ## 检查buf代码风格
	@echo ">> linting buf code..."
	@cd $(PROTO_DIR) && buf lint
	@echo ">> buf code linting done."

buf-breaking: ## 检查buf代码破坏性变更
	@echo ">> checking buf code breaking changes..."
	@cd $(PROTO_DIR) && buf breaking
	@echo ">> buf code breaking changes checking done."

buf-clean: ## 清理生成的buf代码
	@echo ">> cleaning generated protobuf files..."
	@find $(PROTO_DIR) -type f \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -delete 2>/dev/null || true
	@echo ">> protobuf files cleaned."

# wire依赖注入代码生成
wire-install: ## 安装wire工具
	@echo ">> installing wire..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo ">> wire installed: $$(which wire)"

wire: ## 生成wire依赖注入代码
	@echo ">> generating wire code..."
	@cd cmd/arcade-agent && wire
	@echo ">> wire code generation done."

wire-clean: ## 清理wire生成的代码
	@echo ">> cleaning wire generated files..."
	@find . -name "wire_gen.go" -type f -delete
	@echo ">> wire files cleaned."
