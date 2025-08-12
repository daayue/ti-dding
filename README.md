# 钉钉群管理工具 (Ti-Dding)

一个基于 Golang 开发的钉钉群组管理工具，支持批量创建群组、成员管理等操作。

## 功能特性

### 🚀 核心功能
- **批量创建群组**: 支持CSV文件批量导入，自动创建钉钉群组
- **群名重复检测**: 智能检测群名是否已存在，避免重复创建
- **成员管理**: 支持批量添加/移除群成员
- **数据持久化**: 本地文件存储，无需外部数据库依赖

### 📊 数据管理
- CSV文件导入群组信息
- JSON格式本地数据存储
- 支持数据导出和备份
- 群组信息查询和管理

## 技术架构

### 开发语言
- **Golang 1.21+**

### 项目结构
```
ti-dding/
├── cmd/                    # 主程序入口
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── dingtalk/         # 钉钉API客户端
│   ├── models/           # 数据模型
│   ├── services/         # 业务逻辑
│   └── storage/          # 数据存储
├── pkg/                   # 公共包
├── configs/               # 配置文件
├── data/                  # 数据文件
├── scripts/               # 脚本文件
└── docs/                  # 文档
```

### 依赖管理
- 使用 Go modules 管理依赖
- 主要依赖包：
  - `github.com/spf13/cobra` - 命令行框架
  - `github.com/spf13/viper` - 配置管理
  - `encoding/csv` - CSV文件处理
  - `encoding/json` - JSON数据处理

## 快速开始

### 环境要求
- Go 1.21 或更高版本
- 钉钉开发者账号和AppKey/AppSecret

### 安装
```bash
git clone https://github.com/your-username/ti-dding.git
cd ti-dding
go mod tidy
go build -o ti-dding cmd/main.go
```

### 配置
1. 复制配置文件模板
```bash
cp configs/config.yaml.example configs/config.yaml
```

2. 编辑配置文件，填入钉钉应用信息
```yaml
dingtalk:
  app_key: "your_app_key"
  app_secret: "your_app_secret"
  access_token: "your_access_token"
```

### 使用示例

#### 批量创建群组
```bash
# 从CSV文件创建群组
./ti-dding create --file groups.csv

# 查看创建结果
./ti-dding list
```

#### 成员管理
```bash
# 添加成员到所有群组
./ti-dding add-member --user-id "user123" --all-groups

# 添加成员到指定群组
./ti-dding add-member --user-id "user123" --group-id "group123"

# 移除成员
./ti-dding remove-member --user-id "user123" --group-id "group123"
```

## CSV文件格式

### 群组创建CSV格式
```csv
群名称,群描述,群主用户ID,群成员用户ID列表
测试群1,这是一个测试群,user123,"user123,user456,user789"
测试群2,另一个测试群,user456,"user123,user456"
```

## 数据存储

### 群组信息存储
群组信息以JSON格式存储在 `data/groups.json` 文件中：
```json
{
  "groups": [
    {
      "id": "cid123456",
      "name": "测试群1",
      "description": "这是一个测试群",
      "owner_id": "user123",
      "member_count": 3,
      "created_at": "2024-01-01T10:00:00Z",
      "status": "active"
    }
  ]
}
```

## 开发计划

### Phase 1: 基础框架 (Week 1)
- [x] 项目结构搭建
- [x] 配置管理
- [x] 钉钉API客户端基础框架

### Phase 2: 核心功能 (Week 2)
- [ ] 群组创建API集成
- [ ] CSV文件解析
- [ ] 群名重复检测
- [ ] 数据存储实现

### Phase 3: 成员管理 (Week 3)
- [ ] 成员添加功能
- [ ] 成员移除功能
- [ ] 批量操作支持

### Phase 4: 优化完善 (Week 4)
- [ ] 错误处理优化
- [ ] 日志记录
- [ ] 单元测试
- [ ] 文档完善

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue 或联系开发团队。 