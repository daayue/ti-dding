package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"ti-dding/internal/config"
	"ti-dding/internal/models"
)

// Client 钉钉API客户端
type Client struct {
	config      *config.Config
	httpClient  *http.Client
	baseURL     string
	accessToken string
}

// NewClient 创建新的钉钉客户端
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:     cfg.DingTalk.BaseURL,
		accessToken: cfg.GetAccessToken(),
	}
}

// GetAccessToken 获取访问令牌
func (c *Client) GetAccessToken() (string, error) {
	if c.accessToken != "" {
		return c.accessToken, nil
	}

	// 通过AppKey和AppSecret获取访问令牌
	if c.config.DingTalk.AppKey == "" || c.config.DingTalk.AppSecret == "" {
		return "", fmt.Errorf("AppKey和AppSecret不能为空")
	}

	params := url.Values{}
	params.Set("appkey", c.config.DingTalk.AppKey)
	params.Set("appsecret", c.config.DingTalk.AppSecret)

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/gettoken?%s", c.baseURL, params.Encode()))
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		Token   string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Errcode != 0 {
		return "", fmt.Errorf("获取访问令牌失败: %s", result.Errmsg)
	}

	c.accessToken = result.Token
	return result.Token, nil
}

// CreateGroup 创建群组
func (c *Client) CreateGroup(req *models.GroupCreateRequest) (*models.GroupCreateResponse, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 构建钉钉API请求参数
	apiReq := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"owner":       req.OwnerID,
		"useridlist":  req.MemberIDs,
	}

	// 设置群组类型（内部群/外部群）
	if req.IsExternal || req.GroupType == "external" {
		// 外部群设置
		apiReq["conversation_type"] = 2 // 2表示外部群
		apiReq["show_history_type"] = 1 // 允许查看历史消息
		apiReq["validation_type"] = 1   // 允许邀请成员
	} else {
		// 内部群设置
		apiReq["conversation_type"] = 1 // 1表示内部群

		// 添加群组设置
		if c.config.Group.DefaultSettings.AllowMemberInvite {
			apiReq["show_history_type"] = 1
		}
		if c.config.Group.DefaultSettings.AllowMemberView {
			apiReq["validation_type"] = 1
		}
	}

	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 发送请求
	url := fmt.Sprintf("%s/chat/create?access_token=%s", c.baseURL, token)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建群组请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		ChatID  string `json:"chatid"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Errcode != 0 {
		return &models.GroupCreateResponse{
			Success: false,
			Message: fmt.Sprintf("创建群组失败: %s", result.Errmsg),
		}, nil
	}

	return &models.GroupCreateResponse{
		GroupID: result.ChatID,
		Success: true,
		Message: "群组创建成功",
	}, nil
}

// GetGroupList 获取群组列表
func (c *Client) GetGroupList() ([]models.Group, error) {
	_, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 获取部门列表（这里简化处理，实际可能需要遍历部门获取群组）
	// 钉钉API没有直接获取所有群组的接口，需要通过其他方式
	// 这里返回空列表，实际实现可能需要结合其他API
	return []models.Group{}, nil
}

// AddGroupMembers 添加群组成员
func (c *Client) AddGroupMembers(groupID string, userIDs []string) error {
	token, err := c.GetAccessToken()
	if err != nil {
		return err
	}

	apiReq := map[string]interface{}{
		"chatid":     groupID,
		"useridlist": userIDs,
	}

	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	url := fmt.Sprintf("%s/chat/addmember?access_token=%s", c.baseURL, token)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("添加成员请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Errcode != 0 {
		return fmt.Errorf("添加成员失败: %s", result.Errmsg)
	}

	return nil
}

// RemoveGroupMembers 移除群组成员
func (c *Client) RemoveGroupMembers(groupID string, userIDs []string) error {
	token, err := c.GetAccessToken()
	if err != nil {
		return err
	}

	apiReq := map[string]interface{}{
		"chatid":     groupID,
		"useridlist": userIDs,
	}

	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	url := fmt.Sprintf("%s/chat/removemember?access_token=%s", c.baseURL, token)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("移除成员请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Errcode != 0 {
		return fmt.Errorf("移除成员失败: %s", result.Errmsg)
	}

	return nil
}

// CheckGroupExists 检查群组是否存在
func (c *Client) CheckGroupExists(groupName string) (bool, error) {
	// 钉钉API没有直接检查群名是否存在的接口
	// 这里可以通过获取群组列表来检查
	// 实际实现可能需要结合其他方式
	_ = groupName // 暂时未使用，避免linter警告
	return false, nil
}
