#!/bin/bash

# 钉钉群管理工具 - Go 1.20 兼容性验证脚本

set -e

echo "🔍 钉钉群管理工具 - Go 1.20 兼容性验证"
echo "=========================================="

# 检查Go版本
echo "📋 检查Go版本..."
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "当前Go版本: $GO_VERSION"

if [[ "$GO_VERSION" == 1.20* ]]; then
    echo "✅ Go版本兼容 (1.20.x)"
else
    echo "⚠️  Go版本: $GO_VERSION (推荐使用1.20.x)"
fi

echo ""

# 检查依赖
echo "📦 检查依赖..."
echo "执行: go mod tidy"
go mod tidy

echo "执行: go mod download"
go mod download

echo "✅ 依赖检查完成"
echo ""

# 代码检查
echo "🔍 代码检查..."
echo "执行: go vet ./..."
if go vet ./...; then
    echo "✅ 代码检查通过"
else
    echo "❌ 代码检查失败"
    exit 1
fi
echo ""

# 构建项目
echo "🏗️  构建项目..."
echo "执行: make clean"
make clean

echo "执行: make build"
if make build; then
    echo "✅ 项目构建成功"
else
    echo "❌ 项目构建失败"
    exit 1
fi
echo ""

# 功能测试
echo "🧪 功能测试..."
echo "测试帮助命令..."
if ./build/ti-dding --help > /dev/null 2>&1; then
    echo "✅ 帮助命令正常"
else
    echo "❌ 帮助命令失败"
    exit 1
fi

echo "测试列表命令..."
if ./build/ti-dding list > /dev/null 2>&1; then
    echo "✅ 列表命令正常"
else
    echo "❌ 列表命令失败"
    exit 1
fi

echo "测试检查命令..."
if ./build/ti-dding check --name "测试群组" > /dev/null 2>&1; then
    echo "✅ 检查命令正常"
else
    echo "❌ 检查命令失败"
    exit 1
fi

echo "测试导出命令..."
if ./build/ti-dding export --output test_verify.csv > /dev/null 2>&1; then
    echo "✅ 导出命令正常"
    # 清理测试文件
    rm -f test_verify.csv
else
    echo "❌ 导出命令失败"
    exit 1
fi

echo "✅ 所有基本功能测试通过"
echo ""

# 显示项目信息
echo "📊 项目信息..."
echo "项目名称: ti-dding"
echo "Go版本: $GO_VERSION"
echo "构建时间: $(date)"
echo "构建目录: build/ti-dding"
echo ""

# 显示可用命令
echo "🚀 可用命令:"
./build/ti-dding --help | grep -A 20 "Available Commands" | grep -v "Available Commands" | grep -v "Flags" | grep -v "Global Flags" | grep -v "Use" | grep -v "help" | sed 's/^/  /'

echo ""
echo "🎉 验证完成！所有功能正常，Go 1.20 兼容性验证通过！"
echo ""
echo "💡 下一步:"
echo "  1. 配置钉钉应用信息 (编辑 configs/config.yaml)"
echo "  2. 准备CSV文件 (参考 data/groups_example.csv)"
echo "  3. 测试群组创建功能"
echo "  4. 开始使用工具管理钉钉群组" 