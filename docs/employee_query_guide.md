# 钉钉企业员工信息查询工具使用指南

## 🎯 工具说明

这个工具可以帮助你获取钉钉企业中的所有员工信息，包括：
- 员工ID (userid)
- 姓名
- 手机号
- 部门
- 职位
- 邮箱

## 🚀 使用方法

### 1. 构建工具

```bash
# 进入tools目录
cd tools

# 构建员工查询工具
go build -o employee_query employee_query.go
```

### 2. 配置钉钉应用

确保你的 `configs/config.yaml` 文件包含正确的钉钉应用信息：

```yaml
dingtalk:
  app_key: "your_app_key_here"
  app_secret: "your_app_secret_here"
  corp_id: "your_corp_id_here"
  base_url: "https://oapi.dingtalk.com"
```

### 3. 运行工具

```bash
# 使用默认配置文件
./employee_query

# 指定配置文件
./employee_query -config ../configs/config.yaml

# 指定输出文件
./employee_query -output my_employees.csv

# 同时指定配置和输出文件
./employee_query -config ../configs/config.yaml -output my_employees.csv
```

## 📋 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-config` | 配置文件路径 | 自动查找 |
| `-output` | 输出CSV文件路径 | `employees.csv` |

## 🔍 功能特性

### 1. 自动获取访问令牌
- 支持AppKey/AppSecret方式
- 支持直接AccessToken方式
- 自动处理令牌过期

### 2. 部门信息获取
- 获取所有部门列表
- 显示部门ID和名称
- 按部门组织员工信息

### 3. 员工信息获取
- 遍历所有部门获取员工
- 获取员工详细信息
- 错误处理和重试机制

### 4. CSV导出
- 标准CSV格式
- 包含所有员工字段
- 支持自定义文件名

## 📊 输出示例

### 控制台输出
```
🔍 钉钉企业员工信息查询工具
==============================
✅ 访问令牌获取成功
📋 企业部门列表 (共 5 个部门):
  - 技术部 (ID: 1)
  - 产品部 (ID: 2)
  - 运营部 (ID: 3)
  - 市场部 (ID: 4)
  - 人事部 (ID: 5)

正在获取部门 '技术部' 的员工信息...
  部门 '技术部' 找到 15 名员工
正在获取部门 '产品部' 的员工信息...
  部门 '产品部' 找到 8 名员工
...

👥 企业员工列表 (共 45 人):
1. 张三 (ID: 123456, 手机: 13800138000, 部门: 技术部)
2. 李四 (ID: 123457, 手机: 13800138001, 部门: 技术部)
...

📁 员工信息已导出到: employees.csv

🎉 查询完成！
```

### CSV文件内容
```csv
员工ID,姓名,手机号,部门,职位,邮箱
123456,张三,13800138000,技术部,高级工程师,zhangsan@company.com
123457,李四,13800138001,技术部,工程师,lisi@company.com
123458,王五,13800138002,产品部,产品经理,wangwu@company.com
...
```

## ⚠️ 注意事项

### 1. 权限要求
确保你的钉钉应用有以下权限：
- **通讯录身份验证**: 获取访问令牌
- **通讯录只读权限**: 读取部门和员工信息
- **群组管理权限**: 创建和管理群组

### 2. 数据限制
- 钉钉API有调用频率限制
- 大量员工时可能需要较长时间
- 建议在非高峰期使用

### 3. 错误处理
- 工具会自动跳过出错的部门
- 会显示详细的错误信息
- 支持部分成功的情况

## 🔧 故障排除

### 常见问题

#### 1. "获取访问令牌失败"
**原因**: AppKey或AppSecret不正确
**解决**: 检查配置文件中的钉钉应用信息

#### 2. "获取部门列表失败"
**原因**: 应用没有通讯录权限
**解决**: 在钉钉开放平台申请相应权限

#### 3. "获取员工详情失败"
**原因**: 单个员工信息获取失败
**解决**: 工具会自动跳过，不影响整体结果

### 调试模式

如果需要更详细的调试信息，可以修改代码添加日志：

```go
// 在相关函数中添加调试信息
fmt.Printf("DEBUG: 请求URL: %s\n", url)
fmt.Printf("DEBUG: 响应状态: %d\n", resp.StatusCode)
```

## 📚 相关资源

- [钉钉开放平台文档](https://open.dingtalk.com/document/)
- [通讯录管理API](https://open.dingtalk.com/document/orgapp/obtain-the-department-list)
- [用户管理API](https://open.dingtalk.com/document/orgapp/obtain-the-user-list)

## 💡 使用建议

### 1. 定期更新
- 建议定期运行工具更新员工信息
- 可以配合群组创建工具使用
- 保持员工信息的准确性

### 2. 数据备份
- 导出的CSV文件建议备份
- 可以用于其他系统导入
- 支持数据分析需求

### 3. 权限管理
- 合理分配应用权限
- 定期审查权限使用情况
- 遵循最小权限原则

---

**工具版本**: 1.0.0  
**最后更新**: 2024年1月  
**兼容性**: Go 1.20+ 