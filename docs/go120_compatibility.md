# Go 1.20 兼容性说明

## 📋 版本要求

**最低Go版本**: 1.20  
**推荐Go版本**: 1.20+  
**测试版本**: 1.20.1

## 🔧 依赖版本调整

为了确保与Go 1.20的兼容性，我们使用了以下依赖版本：

### 核心依赖
```go
require (
    github.com/spf13/cobra v1.7.0    // 兼容Go 1.20
    github.com/spf13/viper v1.16.0   // 兼容Go 1.20
)
```

### 版本选择原因

1. **cobra v1.7.0**: 
   - 完全兼容Go 1.20
   - 功能稳定，bug修复完善
   - 不依赖Go 1.21+的新特性

2. **viper v1.16.0**:
   - 兼容Go 1.20
   - 配置管理功能完整
   - 性能良好

## ✅ 兼容性验证

### 已测试功能
- [x] 项目构建
- [x] 命令行界面
- [x] 配置管理
- [x] 数据模型
- [x] 存储模块
- [x] 服务层
- [x] 钉钉API客户端

### 测试命令
```bash
# 构建项目
go build -o ti-dding cmd/main.go

# 代码检查
go vet ./...

# 依赖管理
go mod tidy
go mod download

# 功能测试
./ti-dding --help
./ti-dding list
./ti-dding check --name "测试群组"
./ti-dding export --output test.csv
```

## 🚨 注意事项

### 1. 依赖版本锁定
- 不要随意升级到更高版本的依赖
- 特别是cobra和viper，需要确保Go 1.20兼容性

### 2. Go版本升级
如果将来升级到Go 1.21+，可以：
```bash
# 更新go.mod
go mod edit -go=1.21

# 升级依赖版本
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest

# 清理并更新
go mod tidy
```

### 3. 功能限制
Go 1.20版本下，某些新特性不可用：
- `slices` 包（Go 1.21+）
- `log/slog` 包（Go 1.21+）
- 某些新的标准库函数

## 🔄 升级路径

### 短期（保持Go 1.20）
- 使用当前依赖版本
- 确保功能稳定性
- 修复已知问题

### 中期（考虑升级）
- 评估Go 1.21+的新特性
- 测试依赖升级的兼容性
- 制定升级计划

### 长期（完全升级）
- 升级到最新的Go版本
- 利用新特性优化代码
- 重构可能过时的代码

## 📚 相关资源

- [Go 1.20 Release Notes](https://golang.org/doc/go1.20)
- [cobra v1.7.0 Documentation](https://github.com/spf13/cobra/tree/v1.7.0)
- [viper v1.16.0 Documentation](https://github.com/spf13/viper/tree/v1.16.0)

## 🧪 测试建议

### 开发环境
```bash
# 确保使用正确的Go版本
go version

# 清理并重新构建
make clean
make build

# 运行所有测试
make test
```

### 生产环境
- 在生产环境中使用相同的Go版本
- 确保依赖版本一致
- 进行完整的功能测试

## 📞 问题反馈

如果遇到Go 1.20兼容性问题，请：
1. 检查Go版本：`go version`
2. 检查依赖版本：`go list -m all`
3. 查看错误日志
4. 提交Issue到项目仓库

---

**最后更新**: 2024年1月  
**Go版本**: 1.20.1  
**状态**: ✅ 完全兼容 