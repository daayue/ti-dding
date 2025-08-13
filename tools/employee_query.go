package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"ti-dding/internal/config"
)

// Department 部门信息
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Employee 员工信息
type Employee struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	Department string `json:"department"`
	Position   string `json:"position"`
	Email      string `json:"email"`
}

func main() {
	// 命令行参数
	configFile := flag.String("config", "", "配置文件路径")
	outputFile := flag.String("output", "employees.csv", "输出CSV文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 钉钉企业员工信息查询工具")
	fmt.Println("==============================")

	// 获取访问令牌
	token, err := getAccessToken(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取访问令牌失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 访问令牌获取成功")

	// 获取部门列表
	depts, err := getDepartmentList(cfg, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取部门列表失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 企业部门列表 (共 %d 个部门):\n", len(depts))
	for _, dept := range depts {
		fmt.Printf("  - %s (ID: %s)\n", dept.Name, dept.ID)
	}

	// 获取员工列表
	employees, err := getEmployeeList(cfg, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取员工列表失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n👥 企业员工列表 (共 %d 人):\n", len(employees))
	for i, emp := range employees {
		fmt.Printf("%d. %s (ID: %s, 手机: %s, 部门: %s)\n",
			i+1, emp.Name, emp.UserID, emp.Mobile, emp.Department)
	}

	// 导出员工信息到CSV
	if err := exportEmployeesToCSV(employees, *outputFile); err != nil {
		fmt.Printf("⚠️  导出员工信息失败: %v\n", err)
	} else {
		fmt.Printf("\n📁 员工信息已导出到: %s\n", *outputFile)
	}

	fmt.Println("\n🎉 查询完成！")
}

// getAccessToken 获取访问令牌
func getAccessToken(cfg *config.Config) (string, error) {
	if cfg.DingTalk.AccessToken != "" {
		return cfg.DingTalk.AccessToken, nil
	}

	if cfg.DingTalk.AppKey == "" || cfg.DingTalk.AppSecret == "" {
		return "", fmt.Errorf("AppKey和AppSecret不能为空")
	}

	params := url.Values{}
	params.Set("appkey", cfg.DingTalk.AppKey)
	params.Set("appsecret", cfg.DingTalk.AppSecret)

	resp, err := http.Get(fmt.Sprintf("%s/gettoken?%s", cfg.DingTalk.BaseURL, params.Encode()))
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

	return result.Token, nil
}

// getDepartmentList 获取部门列表
func getDepartmentList(cfg *config.Config, token string) ([]Department, error) {
	url := fmt.Sprintf("%s/department/list?access_token=%s", cfg.DingTalk.BaseURL, token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取部门列表失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Errcode int          `json:"errcode"`
		Errmsg  string       `json:"errmsg"`
		Depts   []Department `json:"department"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析部门列表失败: %w", err)
	}

	if result.Errcode != 0 {
		return nil, fmt.Errorf("获取部门列表失败: %s", result.Errmsg)
	}

	return result.Depts, nil
}

// getEmployeeList 获取员工列表
func getEmployeeList(cfg *config.Config, token string) ([]Employee, error) {
	var allEmployees []Employee

	// 获取部门列表
	depts, err := getDepartmentList(cfg, token)
	if err != nil {
		return nil, err
	}

	// 遍历每个部门获取员工
	for _, dept := range depts {
		fmt.Printf("正在获取部门 '%s' 的员工信息...\n", dept.Name)

		url := fmt.Sprintf("%s/user/simplelist?access_token=%s&department_id=%s",
			cfg.DingTalk.BaseURL, token, dept.ID)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("⚠️  获取部门 %s 员工失败: %v\n", dept.Name, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var result struct {
			Errcode int    `json:"errcode"`
			Errmsg  string `json:"errmsg"`
			Users   []struct {
				UserID string `json:"userid"`
				Name   string `json:"name"`
			} `json:"userlist"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}

		if result.Errcode != 0 {
			fmt.Printf("⚠️  部门 %s 返回错误: %s\n", dept.Name, result.Errmsg)
			continue
		}

		fmt.Printf("  部门 '%s' 找到 %d 名员工\n", dept.Name, len(result.Users))

		// 获取员工详细信息
		for _, user := range result.Users {
			emp, err := getEmployeeDetail(cfg, token, user.UserID)
			if err != nil {
				fmt.Printf("⚠️  获取员工 %s 详情失败: %v\n", user.Name, err)
				continue
			}
			emp.Department = dept.Name
			allEmployees = append(allEmployees, emp)
		}
	}

	return allEmployees, nil
}

// getEmployeeDetail 获取员工详细信息
func getEmployeeDetail(cfg *config.Config, token, userID string) (Employee, error) {
	url := fmt.Sprintf("%s/user/get?access_token=%s&userid=%s",
		cfg.DingTalk.BaseURL, token, userID)

	resp, err := http.Get(url)
	if err != nil {
		return Employee{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Employee{}, err
	}

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		User    struct {
			UserID   string `json:"userid"`
			Name     string `json:"name"`
			Mobile   string `json:"mobile"`
			Position string `json:"position"`
			Email    string `json:"email"`
		} `json:"userinfo"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Employee{}, err
	}

	if result.Errcode != 0 {
		return Employee{}, fmt.Errorf("获取员工详情失败: %s", result.Errmsg)
	}

	return Employee{
		UserID:   result.User.UserID,
		Name:     result.User.Name,
		Mobile:   result.User.Mobile,
		Position: result.User.Position,
		Email:    result.User.Email,
	}, nil
}

// exportEmployeesToCSV 导出员工信息到CSV
func exportEmployeesToCSV(employees []Employee, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入标题行
	header := "员工ID,姓名,手机号,部门,职位,邮箱\n"
	if _, err := file.WriteString(header); err != nil {
		return err
	}

	// 写入数据行
	for _, emp := range employees {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s\n",
			emp.UserID, emp.Name, emp.Mobile, emp.Department, emp.Position, emp.Email)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}
