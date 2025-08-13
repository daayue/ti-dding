package services

import (
	"fmt"
	"strings"

	"ti-dding/internal/config"
	"ti-dding/internal/dingtalk"
	"ti-dding/internal/models"
	"ti-dding/internal/storage"
)

// GroupService 群组服务
type GroupService struct {
	dingtalkClient *dingtalk.Client
	storage        storage.Storage
	config         *config.GroupConfig
}

// NewGroupService 创建新的群组服务
func NewGroupService(client *dingtalk.Client, storage storage.Storage, config *config.GroupConfig) *GroupService {
	return &GroupService{
		dingtalkClient: client,
		storage:        storage,
		config:         config,
	}
}

// CreateGroupsFromCSV 从CSV文件批量创建群组
func (s *GroupService) CreateGroupsFromCSV(csvFile string) (*models.GroupCreateResponse, error) {
	// 从CSV文件加载群组数据
	csvGroups, err := s.storage.(*storage.FileStorage).LoadGroupsFromCSV(csvFile)
	if err != nil {
		return nil, fmt.Errorf("加载CSV文件失败: %w", err)
	}

	if len(csvGroups) == 0 {
		return &models.GroupCreateResponse{
			Success: false,
			Message: "CSV文件中没有有效的群组数据",
		}, nil
	}

	var successCount, failCount int
	var failedGroups []string

	// 逐个创建群组
	for _, csvGroup := range csvGroups {
		// 检查群名是否已存在
		if s.storage.GroupExists(csvGroup.Name) {
			failedGroups = append(failedGroups, fmt.Sprintf("%s (群名已存在)", csvGroup.Name))
			failCount++
			continue
		}

		// 解析成员ID列表
		memberIDs := []string{}
		if csvGroup.MemberIDs != "" {
			memberIDs = strings.Split(csvGroup.MemberIDs, ",")
			// 清理空白字符
			for i, id := range memberIDs {
				memberIDs[i] = strings.TrimSpace(id)
			}
		}

		// 确保群主在成员列表中
		ownerInMembers := false
		for _, memberID := range memberIDs {
			if memberID == csvGroup.OwnerID {
				ownerInMembers = true
				break
			}
		}
		if !ownerInMembers {
			memberIDs = append(memberIDs, csvGroup.OwnerID)
		}

		// 确定群组类型
		groupType := "internal"
		isExternal := false
		if csvGroup.GroupType != "" {
			switch strings.ToLower(strings.TrimSpace(csvGroup.GroupType)) {
			case "external", "外部群", "外部":
				groupType = "external"
				isExternal = true
			case "internal", "内部群", "内部":
				groupType = "internal"
				isExternal = false
			default:
				groupType = "internal"
				isExternal = false
			}
		}

		// 创建群组请求
		req := &models.GroupCreateRequest{
			Name:        csvGroup.Name,
			Description: csvGroup.Description,
			OwnerID:     csvGroup.OwnerID,
			MemberIDs:   memberIDs,
			GroupType:   groupType,
			IsExternal:  isExternal,
		}

		// 调用钉钉API创建群组
		resp, err := s.dingtalkClient.CreateGroup(req)
		if err != nil {
			failedGroups = append(failedGroups, fmt.Sprintf("%s (API调用失败: %s)", csvGroup.Name, err.Error()))
			failCount++
			continue
		}

		if !resp.Success {
			failedGroups = append(failedGroups, fmt.Sprintf("%s (%s)", csvGroup.Name, resp.Message))
			failCount++
			continue
		}

		// 创建成功，保存到本地存储
		group := models.NewGroup(csvGroup.Name, csvGroup.Description, csvGroup.OwnerID)
		group.ID = resp.GroupID
		group.Members = memberIDs
		group.MemberCount = len(memberIDs)

		if err := s.storage.AddGroup(*group); err != nil {
			failedGroups = append(failedGroups, fmt.Sprintf("%s (保存失败: %s)", csvGroup.Name, err.Error()))
			failCount++
			continue
		}

		successCount++
	}

	// 构建响应消息
	var message string
	if successCount > 0 {
		message = fmt.Sprintf("成功创建 %d 个群组", successCount)
		if failCount > 0 {
			message += fmt.Sprintf("，失败 %d 个群组", failCount)
		}
	} else {
		message = "没有成功创建任何群组"
	}

	if len(failedGroups) > 0 {
		message += "\n失败的群组：" + strings.Join(failedGroups, "; ")
	}

	return &models.GroupCreateResponse{
		Success: successCount > 0,
		Message: message,
	}, nil
}

// ListGroups 获取群组列表
func (s *GroupService) ListGroups() (*models.GroupListResponse, error) {
	groups, err := s.storage.LoadGroups()
	if err != nil {
		return nil, fmt.Errorf("加载群组列表失败: %w", err)
	}

	// 过滤掉已删除的群组
	var activeGroups []models.Group
	for _, group := range groups {
		if group.Status != "deleted" {
			activeGroups = append(activeGroups, group)
		}
	}

	return &models.GroupListResponse{
		Groups: activeGroups,
		Total:  len(activeGroups),
	}, nil
}

// AddMembers 添加成员到群组
func (s *GroupService) AddMembers(req *models.GroupMemberRequest) (*models.GroupMemberResponse, error) {
	if len(req.UserIDs) == 0 {
		return &models.GroupMemberResponse{
			Success: false,
			Message: "用户ID列表不能为空",
		}, nil
	}

	var affectedGroups int
	var errors []string

	if req.AllGroups {
		// 添加到所有群组
		groups, err := s.storage.LoadGroups()
		if err != nil {
			return nil, fmt.Errorf("加载群组列表失败: %w", err)
		}

		for _, group := range groups {
			if group.Status == "deleted" {
				continue
			}

			// 调用钉钉API添加成员
			if err := s.dingtalkClient.AddGroupMembers(group.ID, req.UserIDs); err != nil {
				errors = append(errors, fmt.Sprintf("群组 %s: %s", group.Name, err.Error()))
				continue
			}

			// 更新本地存储
			for _, userID := range req.UserIDs {
				group.AddMember(userID)
			}
			if err := s.storage.UpdateGroup(group); err != nil {
				errors = append(errors, fmt.Sprintf("群组 %s 更新失败: %s", group.Name, err.Error()))
				continue
			}

			affectedGroups++
		}
	} else {
		// 添加到指定群组
		if req.GroupID == "" {
			return &models.GroupMemberResponse{
				Success: false,
				Message: "指定群组操作时必须提供群组ID",
			}, nil
		}

		group, err := s.storage.GetGroupByID(req.GroupID)
		if err != nil {
			return &models.GroupMemberResponse{
				Success: false,
				Message: fmt.Sprintf("群组不存在: %s", err.Error()),
			}, nil
		}

		// 调用钉钉API添加成员
		if err := s.dingtalkClient.AddGroupMembers(group.ID, req.UserIDs); err != nil {
			return &models.GroupMemberResponse{
				Success: false,
				Message: fmt.Sprintf("添加成员失败: %s", err.Error()),
			}, nil
		}

		// 更新本地存储
		for _, userID := range req.UserIDs {
			group.AddMember(userID)
		}
		if err := s.storage.UpdateGroup(*group); err != nil {
			return &models.GroupMemberResponse{
				Success: false,
				Message: fmt.Sprintf("更新群组信息失败: %s", err.Error()),
			}, nil
		}

		affectedGroups = 1
	}

	// 构建响应消息
	var message string
	if len(errors) == 0 {
		message = fmt.Sprintf("成功添加成员到 %d 个群组", affectedGroups)
	} else {
		message = fmt.Sprintf("部分成功：%d 个群组，错误：%s", affectedGroups, strings.Join(errors, "; "))
	}

	return &models.GroupMemberResponse{
		Success:  affectedGroups > 0,
		Message:  message,
		Affected: affectedGroups,
	}, nil
}

// RemoveMembers 从群组移除成员
func (s *GroupService) RemoveMembers(req *models.GroupMemberRequest) (*models.GroupMemberResponse, error) {
	if len(req.UserIDs) == 0 {
		return &models.GroupMemberResponse{
			Success: false,
			Message: "用户ID列表不能为空",
		}, nil
	}

	var affectedGroups int
	var errors []string

	if req.AllGroups {
		// 从所有群组移除
		groups, err := s.storage.LoadGroups()
		if err != nil {
			return nil, fmt.Errorf("加载群组列表失败: %w", err)
		}

		for _, group := range groups {
			if group.Status == "deleted" {
				continue
			}

			// 调用钉钉API移除成员
			if err := s.dingtalkClient.RemoveGroupMembers(group.ID, req.UserIDs); err != nil {
				errors = append(errors, fmt.Sprintf("群组 %s: %s", group.Name, err.Error()))
				continue
			}

			// 更新本地存储
			for _, userID := range req.UserIDs {
				group.RemoveMember(userID)
			}
			if err := s.storage.UpdateGroup(group); err != nil {
				errors = append(errors, fmt.Sprintf("群组 %s 更新失败: %s", group.Name, err.Error()))
				continue
			}

			affectedGroups++
		}
	} else {
		// 从指定群组移除
		if req.GroupID == "" {
			return &models.GroupMemberResponse{
				Success: false,
				Message: "指定群组操作时必须提供群组ID",
			}, nil
		}

		group, err := s.storage.GetGroupByID(req.GroupID)
		if err != nil {
			return &models.GroupMemberResponse{
				Success: false,
				Message: fmt.Sprintf("群组不存在: %s", err.Error()),
			}, nil
		}

		// 调用钉钉API移除成员
		if err := s.dingtalkClient.RemoveGroupMembers(group.ID, req.UserIDs); err != nil {
			return &models.GroupMemberResponse{
				Success: false,
				Message: fmt.Sprintf("移除成员失败: %s", err.Error()),
			}, nil
		}

		// 更新本地存储
		for _, userID := range req.UserIDs {
			group.RemoveMember(userID)
		}
		if err := s.storage.UpdateGroup(*group); err != nil {
			return &models.GroupMemberResponse{
				Success: false,
				Message: fmt.Sprintf("更新群组信息失败: %s", err.Error()),
			}, nil
		}

		affectedGroups = 1
	}

	// 构建响应消息
	var message string
	if len(errors) == 0 {
		message = fmt.Sprintf("成功从 %d 个群组移除成员", affectedGroups)
	} else {
		message = fmt.Sprintf("部分成功：%d 个群组，错误：%s", affectedGroups, strings.Join(errors, "; "))
	}

	return &models.GroupMemberResponse{
		Success:  affectedGroups > 0,
		Message:  message,
		Affected: affectedGroups,
	}, nil
}

// CheckGroupExists 检查群组是否存在
func (s *GroupService) CheckGroupExists(name string) bool {
	return s.storage.GroupExists(name)
}

// ExportGroups 导出群组数据
func (s *GroupService) ExportGroups(outputFile string) error {
	return s.storage.(*storage.FileStorage).ExportGroupsToCSV(outputFile)
}
