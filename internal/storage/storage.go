package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ti-dding/internal/models"
)

// Storage 数据存储接口
type Storage interface {
	SaveGroups(groups []models.Group) error
	LoadGroups() ([]models.Group, error)
	AddGroup(group models.Group) error
	UpdateGroup(group models.Group) error
	DeleteGroup(groupID string) error
	GetGroupByID(groupID string) (*models.Group, error)
	GetGroupByName(name string) (*models.Group, error)
	GroupExists(name string) bool
}

// FileStorage 文件存储实现
type FileStorage struct {
	dataDir    string
	groupsFile string
}

// NewFileStorage 创建新的文件存储实例
func NewFileStorage(dataDir string) *FileStorage {
	return &FileStorage{
		dataDir:    dataDir,
		groupsFile: filepath.Join(dataDir, "groups.json"),
	}
}

// SaveGroups 保存群组列表到文件
func (fs *FileStorage) SaveGroups(groups []models.Group) error {
	// 确保数据目录存在
	if err := os.MkdirAll(fs.dataDir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	// 准备数据
	data := struct {
		Groups    []models.Group `json:"groups"`
		Total     int            `json:"total"`
		UpdatedAt time.Time      `json:"updated_at"`
	}{
		Groups:    groups,
		Total:     len(groups),
		UpdatedAt: time.Now(),
	}

	// 序列化为JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化群组数据失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(fs.groupsFile, jsonData, 0644); err != nil {
		return fmt.Errorf("写入群组数据文件失败: %w", err)
	}

	return nil
}

// LoadGroups 从文件加载群组列表
func (fs *FileStorage) LoadGroups() ([]models.Group, error) {
	// 检查文件是否存在
	if _, err := os.Stat(fs.groupsFile); os.IsNotExist(err) {
		// 文件不存在，返回空列表
		return []models.Group{}, nil
	}

	// 读取文件内容
	jsonData, err := os.ReadFile(fs.groupsFile)
	if err != nil {
		return nil, fmt.Errorf("读取群组数据文件失败: %w", err)
	}

	// 解析JSON数据
	var data struct {
		Groups []models.Group `json:"groups"`
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("解析群组数据失败: %w", err)
	}

	return data.Groups, nil
}

// AddGroup 添加新群组
func (fs *FileStorage) AddGroup(group models.Group) error {
	groups, err := fs.LoadGroups()
	if err != nil {
		return err
	}

	// 检查群组是否已存在
	for _, existingGroup := range groups {
		if existingGroup.ID == group.ID || existingGroup.Name == group.Name {
			return fmt.Errorf("群组已存在: ID=%s, Name=%s", group.ID, group.Name)
		}
	}

	// 添加新群组
	groups = append(groups, group)

	// 保存到文件
	return fs.SaveGroups(groups)
}

// UpdateGroup 更新群组信息
func (fs *FileStorage) UpdateGroup(group models.Group) error {
	groups, err := fs.LoadGroups()
	if err != nil {
		return err
	}

	// 查找并更新群组
	for i, existingGroup := range groups {
		if existingGroup.ID == group.ID {
			groups[i] = group
			return fs.SaveGroups(groups)
		}
	}

	return fmt.Errorf("群组不存在: ID=%s", group.ID)
}

// DeleteGroup 删除群组
func (fs *FileStorage) DeleteGroup(groupID string) error {
	groups, err := fs.LoadGroups()
	if err != nil {
		return err
	}

	// 查找并删除群组
	for i, group := range groups {
		if group.ID == groupID {
			// 标记为已删除而不是物理删除
			groups[i].Status = "deleted"
			groups[i].UpdatedAt = time.Now()
			return fs.SaveGroups(groups)
		}
	}

	return fmt.Errorf("群组不存在: ID=%s", groupID)
}

// GetGroupByID 根据ID获取群组
func (fs *FileStorage) GetGroupByID(groupID string) (*models.Group, error) {
	groups, err := fs.LoadGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.ID == groupID && group.Status != "deleted" {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("群组不存在: ID=%s", groupID)
}

// GetGroupByName 根据名称获取群组
func (fs *FileStorage) GetGroupByName(name string) (*models.Group, error) {
	groups, err := fs.LoadGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Name == name && group.Status != "deleted" {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("群组不存在: Name=%s", name)
}

// GroupExists 检查群组是否存在
func (fs *FileStorage) GroupExists(name string) bool {
	groups, err := fs.LoadGroups()
	if err != nil {
		return false
	}

	for _, group := range groups {
		if group.Name == name && group.Status != "deleted" {
			return true
		}
	}

	return false
}

// LoadGroupsFromCSV 从CSV文件加载群组数据
func (fs *FileStorage) LoadGroupsFromCSV(csvFile string) ([]models.CSVGroupData, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("打开CSV文件失败: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // 允许变长记录

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV文件失败: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV文件格式错误：至少需要标题行和一行数据")
	}

	var groups []models.CSVGroupData

	// 跳过标题行，从第二行开始
	for i, record := range records[1:] {
		if len(record) < 4 {
			return nil, fmt.Errorf("第%d行数据不完整，需要至少4个字段", i+2)
		}

		// 处理群组类型字段（可选）
		groupType := ""
		if len(record) > 4 {
			groupType = strings.TrimSpace(record[4])
		}

		group := models.CSVGroupData{
			Name:        strings.TrimSpace(record[0]),
			Description: strings.TrimSpace(record[1]),
			OwnerID:     strings.TrimSpace(record[2]),
			MemberIDs:   strings.TrimSpace(record[3]),
			GroupType:   groupType,
		}

		// 验证必填字段
		if group.Name == "" {
			return nil, fmt.Errorf("第%d行群名称不能为空", i+2)
		}
		if group.OwnerID == "" {
			return nil, fmt.Errorf("第%d行群主用户ID不能为空", i+2)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// ExportGroupsToCSV 导出群组数据到CSV文件
func (fs *FileStorage) ExportGroupsToCSV(outputFile string) error {
	groups, err := fs.LoadGroups()
	if err != nil {
		return err
	}

	// 过滤掉已删除的群组
	var activeGroups []models.Group
	for _, group := range groups {
		if group.Status != "deleted" {
			activeGroups = append(activeGroups, group)
		}
	}

	// 创建CSV文件
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建CSV文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入标题行
	headers := []string{"群组ID", "群名称", "群描述", "群主用户ID", "成员数量", "群组类型", "创建时间", "状态"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("写入CSV标题失败: %w", err)
	}

	// 写入数据行
	for _, group := range activeGroups {
		// 确定群组类型显示文本
		groupTypeText := "内部群"
		if group.IsExternal {
			groupTypeText = "外部群"
		}

		record := []string{
			group.ID,
			group.Name,
			group.Description,
			group.OwnerID,
			fmt.Sprintf("%d", group.MemberCount),
			groupTypeText,
			group.CreatedAt.Format("2006-01-02 15:04:05"),
			group.Status,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	return nil
}
