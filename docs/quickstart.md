# 钉钉群管理工具 - 快速开始指南

## 🚀 快速开始

### 1. 环境准备

确保你的系统已安装：
- Go 1.20 或更高版本
- Git

### 2. 获取代码

```bash
git clone https://github.com/your-username/ti-dding.git
cd ti-dding
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置钉钉应用

1. 复制配置文件模板：
```bash
cp configs/config.yaml.example configs/config.yaml
```

2. 编辑配置文件，填入你的钉钉应用信息：
```yaml
dingtalk:
  app_key: "your_app_key_here"        # 从钉钉开放平台获取
  app_secret: "your_app_secret_here"  # 从钉钉开放平台获取
  corp_id: "your_corp_id_here"        # 企业ID（可选）
```

### 5. 构建应用

```bash
make build
```

或者手动构建：
```bash
go build -o ti-dding cmd/main.go
```

### 6. 测试工具

查看帮助信息：
```bash
./ti-dding --help
```

## 📋 使用示例

### 批量创建群组

1. 准备CSV文件（参考 `data/groups_example.csv`）：
```csv
群名称,群描述,群主用户ID,群成员用户ID列表
测试群1,这是一个测试群组,user123,"user123,user456,user789"
测试群2,另一个测试群组,user456,"user456,user123"
```

2. 执行创建命令：
```bash
./ti-dding create --file groups.csv
```

### 查看群组列表

```bash
./ti-dding list
```

### 添加成员

添加到指定群组：
```bash
./ti-dding add-member --user-id "user101" --group-id "group123"
```

添加到所有群组：
```bash
./ti-dding add-member --user-id "user101" --all-groups
```

### 移除成员

从指定群组移除：
```bash
./ti-dding remove-member --user-id "user101" --group-id "group123"
```

从所有群组移除：
```bash
./ti-dding remove-member --user-id "user101" --all-groups
```

### 检查群组是否存在

```bash
./ti-dding check --name "测试群1"
```

### 导出群组数据

```bash
./ti-dding export --output "groups_export.csv"
```

## 🔧 常用命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `create` | 批量创建群组 | `./ti-dding create --file groups.csv` |
| `list` | 列出所有群组 | `./ti-dding list` |
| `add-member` | 添加群组成员 | `./ti-dding add-member --user-id user123 --all-groups` |
| `remove-member` | 移除群组成员 | `./ti-dding remove-member --user-id user123 --group-id group123` |
| `check` | 检查群组是否存在 | `./ti-dding check --name "群组名称"` |
| `export` | 导出群组数据 | `./ti-dding export --output export.csv` |

## 📁 文件结构

```
ti-dding/
├── cmd/main.go              # 主程序入口
├── internal/                # 内部包
│   ├── config/             # 配置管理
│   ├── dingtalk/           # 钉钉API客户端
│   ├── models/             # 数据模型
│   ├── services/           # 业务逻辑
│   └── storage/            # 数据存储
├── configs/                 # 配置文件
├── data/                    # 数据文件
├── docs/                    # 文档
├── Makefile                 # 构建脚本
└── README.md               # 项目说明
```

## ⚠️ 注意事项

1. **权限要求**：确保你的钉钉应用有创建群组和管理成员的权限
2. **用户ID格式**：用户ID必须是钉钉系统中的有效用户ID
3. **群组名称唯一性**：群组名称在系统中必须唯一
4. **数据备份**：定期备份 `data/` 目录中的数据文件

## 🆘 常见问题

### Q: 如何获取钉钉应用的AppKey和AppSecret？
A: 登录钉钉开放平台，在应用管理页面可以找到应用的AppKey和AppSecret。

### Q: 为什么创建群组失败？
A: 可能的原因：
- 配置文件中的AppKey/AppSecret不正确
- 应用没有创建群组的权限
- 群组名称已存在
- 用户ID不存在或无效

### Q: 如何批量添加多个用户？
A: 目前支持一次添加一个用户，如需批量添加多个用户，可以多次执行命令或修改代码支持批量操作。

### Q: 数据存储在哪里？
A: 群组数据存储在 `data/groups.json` 文件中，路径可在配置文件中自定义。

## 📞 获取帮助

如果遇到问题，可以：
1. 查看帮助信息：`./ti-dding --help`
2. 查看具体命令帮助：`./ti-dding create --help`
3. 提交Issue到项目仓库
4. 查看项目README文档

## 🔄 更新工具

```bash
git pull origin main
go mod tidy
make build
``` 