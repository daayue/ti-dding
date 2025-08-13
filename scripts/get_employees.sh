#!/bin/bash

# 钉钉企业员工信息获取脚本

echo "🔍 钉钉企业员工信息获取工具"
echo "=============================="

# 检查工具是否存在
if [ ! -f "tools/employee_query" ]; then
    echo "❌ 员工查询工具不存在，正在构建..."
    cd tools
    go build -o employee_query employee_query.go
    cd ..
    
    if [ ! -f "tools/employee_query" ]; then
        echo "❌ 构建失败，请检查Go环境"
        exit 1
    fi
    echo "✅ 工具构建成功"
fi

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo "⚠️  配置文件不存在，正在创建..."
    cp configs/config.yaml.example configs/config.yaml
    echo "📝 请编辑 configs/config.yaml 文件，填入你的钉钉应用信息"
    echo "   需要的信息："
    echo "   - app_key: 钉钉应用的AppKey"
    echo "   - app_secret: 钉钉应用的AppSecret"
    echo "   - corp_id: 企业ID（可选）"
    echo ""
    echo "   获取方式："
    echo "   1. 登录钉钉开放平台 (https://open.dingtalk.com/)"
    echo "   2. 进入应用管理 -> 你的应用"
    echo "   3. 查看应用信息页面"
    echo ""
    read -p "配置完成后按回车键继续..."
fi

# 运行员工查询工具
echo "🚀 开始获取员工信息..."
echo ""

# 使用配置文件运行工具
./tools/employee_query -config configs/config.yaml -output employees_$(date +%Y%m%d_%H%M%S).csv

echo ""
echo "📋 操作完成！"
echo ""
echo "💡 下一步操作："
echo "   1. 查看生成的CSV文件，获取正确的员工ID"
echo "   2. 使用员工ID更新你的群组创建CSV文件"
echo "   3. 运行群组创建命令："
echo "      ./ti-dding create --file your_groups.csv"
echo ""
echo "📚 更多帮助请查看："
echo "   - docs/employee_query_guide.md"
echo "   - docs/quickstart.md" 