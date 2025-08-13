package models

import (
	"time"
)

// Group 群组信息
type Group struct {
	ID          string    `json:"id"`           // 群组ID
	Name        string    `json:"name"`         // 群组名称
	Description string    `json:"description"`  // 群组描述
	OwnerID     string    `json:"owner_id"`     // 群主用户ID
	MemberCount int       `json:"member_count"` // 成员数量
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
	Status      string    `json:"status"`       // 群组状态: active, inactive, deleted
	Members     []string  `json:"members"`      // 成员用户ID列表
	GroupType   string    `json:"group_type"`   // 群组类型: internal(内部群), external(外部群)
	IsExternal  bool      `json:"is_external"`  // 是否为外部群
}

// GroupCreateRequest 创建群组请求
type GroupCreateRequest struct {
	Name        string   `json:"name"`        // 群组名称
	Description string   `json:"description"` // 群组描述
	OwnerID     string   `json:"owner_id"`    // 群主用户ID
	MemberIDs   []string `json:"member_ids"`  // 初始成员用户ID列表
	GroupType   string   `json:"group_type"`  // 群组类型: internal(内部群), external(外部群)
	IsExternal  bool     `json:"is_external"` // 是否为外部群
}

// GroupCreateResponse 创建群组响应
type GroupCreateResponse struct {
	GroupID string `json:"group_id"` // 群组ID
	Success bool   `json:"success"`  // 是否成功
	Message string `json:"message"`  // 响应消息
}

// GroupListResponse 群组列表响应
type GroupListResponse struct {
	Groups []Group `json:"groups"` // 群组列表
	Total  int     `json:"total"`  // 总数
}

// GroupMemberRequest 群组成员操作请求
type GroupMemberRequest struct {
	GroupID   string   `json:"group_id"`   // 群组ID
	UserIDs   []string `json:"user_ids"`   // 用户ID列表
	AllGroups bool     `json:"all_groups"` // 是否操作所有群组
}

// GroupMemberResponse 群组成员操作响应
type GroupMemberResponse struct {
	Success  bool   `json:"success"`  // 是否成功
	Message  string `json:"message"`  // 响应消息
	Affected int    `json:"affected"` // 影响的群组数量
}

// CSVGroupData CSV文件中的群组数据
type CSVGroupData struct {
	Name        string `csv:"群名称"`
	Description string `csv:"群描述"`
	OwnerID     string `csv:"群主用户ID"`
	MemberIDs   string `csv:"群成员用户ID列表"`
	GroupType   string `csv:"群组类型"` // 内部群/外部群
}

// NewGroup 创建新的群组实例
func NewGroup(name, description, ownerID string) *Group {
	return NewGroupWithType(name, description, ownerID, "internal", false)
}

// NewGroupWithType 创建指定类型的群组实例
func NewGroupWithType(name, description, ownerID, groupType string, isExternal bool) *Group {
	now := time.Now()
	return &Group{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		MemberCount: 1, // 初始只有群主
		CreatedAt:   now,
		UpdatedAt:   now,
		Status:      "active",
		Members:     []string{ownerID},
		GroupType:   groupType,
		IsExternal:  isExternal,
	}
}

// AddMember 添加成员
func (g *Group) AddMember(userID string) bool {
	// 检查是否已经是成员
	for _, member := range g.Members {
		if member == userID {
			return false
		}
	}

	g.Members = append(g.Members, userID)
	g.MemberCount = len(g.Members)
	g.UpdatedAt = time.Now()
	return true
}

// RemoveMember 移除成员
func (g *Group) RemoveMember(userID string) bool {
	for i, member := range g.Members {
		if member == userID {
			// 不能移除群主
			if member == g.OwnerID {
				return false
			}

			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			g.MemberCount = len(g.Members)
			g.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// IsMember 检查用户是否是群成员
func (g *Group) IsMember(userID string) bool {
	for _, member := range g.Members {
		if member == userID {
			return true
		}
	}
	return false
}

// IsOwner 检查用户是否是群主
func (g *Group) IsOwner(userID string) bool {
	return g.OwnerID == userID
}
