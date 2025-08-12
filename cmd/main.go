package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ti-dding/internal/config"
	"ti-dding/internal/dingtalk"
	"ti-dding/internal/models"
	"ti-dding/internal/services"
	"ti-dding/internal/storage"
)

var (
	configFile string
	cfg        *config.Config
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "ti-dding",
	Short: "钉钉群管理工具",
	Long: `钉钉群管理工具 (Ti-Dding)

一个基于 Golang 开发的钉钉群组管理工具，支持批量创建群组、成员管理等操作。

主要功能：
- 批量创建群组 (支持CSV文件导入)
- 群组成员管理 (添加/移除成员)
- 群组信息查询和管理
- 本地数据存储和导出

使用示例：
  ti-dding create --file groups.csv    # 从CSV文件创建群组
  ti-dding list                       # 查看群组列表
  ti-dding add-member --user-id user123 --all-groups  # 添加成员到所有群组
  ti-dding export --output groups.csv # 导出群组数据`,
}

// createCmd 创建群组命令
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建群组",
	Long:  "从CSV文件批量创建钉钉群组",
	RunE: func(cmd *cobra.Command, args []string) error {
		csvFile, _ := cmd.Flags().GetString("file")
		if csvFile == "" {
			return fmt.Errorf("必须指定CSV文件路径 (--file)")
		}

		// 初始化服务
		client := dingtalk.NewClient(cfg)
		storage := storage.NewFileStorage(cfg.GetDataDir())
		groupConfig := &config.GroupConfig{
			DefaultOwner: cfg.Group.DefaultOwner,
			DefaultSettings: config.GroupDefaultSettings{
				AllowMemberInvite:   cfg.Group.DefaultSettings.AllowMemberInvite,
				AllowMemberView:     cfg.Group.DefaultSettings.AllowMemberView,
				AllowMemberEditName: cfg.Group.DefaultSettings.AllowMemberEditName,
			},
		}
		service := services.NewGroupService(client, storage, groupConfig)

		// 执行创建操作
		resp, err := service.CreateGroupsFromCSV(csvFile)
		if err != nil {
			return fmt.Errorf("创建群组失败: %w", err)
		}

		fmt.Println(resp.Message)
		return nil
	},
}

// listCmd 列出群组命令
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出群组",
	Long:  "显示所有群组的列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化服务
		client := dingtalk.NewClient(cfg)
		storage := storage.NewFileStorage(cfg.GetDataDir())
		groupConfig := &config.GroupConfig{
			DefaultOwner: cfg.Group.DefaultOwner,
			DefaultSettings: config.GroupDefaultSettings{
				AllowMemberInvite:   cfg.Group.DefaultSettings.AllowMemberInvite,
				AllowMemberView:     cfg.Group.DefaultSettings.AllowMemberView,
				AllowMemberEditName: cfg.Group.DefaultSettings.AllowMemberEditName,
			},
		}
		service := services.NewGroupService(client, storage, groupConfig)

		// 获取群组列表
		resp, err := service.ListGroups()
		if err != nil {
			return fmt.Errorf("获取群组列表失败: %w", err)
		}

		if resp.Total == 0 {
			fmt.Println("暂无群组")
			return nil
		}

		fmt.Printf("共有 %d 个群组:\n\n", resp.Total)
		for i, group := range resp.Groups {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, group.Name, group.ID)
			fmt.Printf("   描述: %s\n", group.Description)
			fmt.Printf("   群主: %s\n", group.OwnerID)
			fmt.Printf("   成员数: %d\n", group.MemberCount)
			fmt.Printf("   创建时间: %s\n", group.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("   状态: %s\n\n", group.Status)
		}

		return nil
	},
}

// addMemberCmd 添加成员命令
var addMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "添加群组成员",
	Long:  "添加用户到指定的群组或所有群组",
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")
		groupID, _ := cmd.Flags().GetString("group-id")
		allGroups, _ := cmd.Flags().GetBool("all-groups")

		if userID == "" {
			return fmt.Errorf("必须指定用户ID (--user-id)")
		}

		if !allGroups && groupID == "" {
			return fmt.Errorf("必须指定群组ID (--group-id) 或使用 --all-groups")
		}

		// 初始化服务
		client := dingtalk.NewClient(cfg)
		storage := storage.NewFileStorage(cfg.GetDataDir())
		groupConfig := &config.GroupConfig{
			DefaultOwner: cfg.Group.DefaultOwner,
			DefaultSettings: config.GroupDefaultSettings{
				AllowMemberInvite:   cfg.Group.DefaultSettings.AllowMemberInvite,
				AllowMemberView:     cfg.Group.DefaultSettings.AllowMemberView,
				AllowMemberEditName: cfg.Group.DefaultSettings.AllowMemberEditName,
			},
		}
		service := services.NewGroupService(client, storage, groupConfig)

		// 执行添加成员操作
		req := &models.GroupMemberRequest{
			UserIDs:   []string{userID},
			GroupID:   groupID,
			AllGroups: allGroups,
		}

		resp, err := service.AddMembers(req)
		if err != nil {
			return fmt.Errorf("添加成员失败: %w", err)
		}

		fmt.Println(resp.Message)
		return nil
	},
}

// removeMemberCmd 移除成员命令
var removeMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "移除群组成员",
	Long:  "从指定的群组或所有群组中移除用户",
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")
		groupID, _ := cmd.Flags().GetString("group-id")
		allGroups, _ := cmd.Flags().GetBool("all-groups")

		if userID == "" {
			return fmt.Errorf("必须指定用户ID (--user-id)")
		}

		if !allGroups && groupID == "" {
			return fmt.Errorf("必须指定群组ID (--group-id) 或使用 --all-groups")
		}

		// 初始化服务
		client := dingtalk.NewClient(cfg)
		storage := storage.NewFileStorage(cfg.GetDataDir())
		groupConfig := &config.GroupConfig{
			DefaultOwner: cfg.Group.DefaultOwner,
			DefaultSettings: config.GroupDefaultSettings{
				AllowMemberInvite:   cfg.Group.DefaultSettings.AllowMemberInvite,
				AllowMemberView:     cfg.Group.DefaultSettings.AllowMemberView,
				AllowMemberEditName: cfg.Group.DefaultSettings.AllowMemberEditName,
			},
		}
		service := services.NewGroupService(client, storage, groupConfig)

		// 执行移除成员操作
		req := &models.GroupMemberRequest{
			UserIDs:   []string{userID},
			GroupID:   groupID,
			AllGroups: allGroups,
		}

		resp, err := service.RemoveMembers(req)
		if err != nil {
			return fmt.Errorf("移除成员失败: %w", err)
		}

		fmt.Println(resp.Message)
		return nil
	},
}

// exportCmd 导出群组命令
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出群组数据",
	Long:  "将群组数据导出为CSV文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = "groups_export.csv"
		}

		// 初始化服务
		client := dingtalk.NewClient(cfg)
		storage := storage.NewFileStorage(cfg.GetDataDir())
		groupConfig := &models.GroupConfig{
			DefaultOwner: cfg.Group.DefaultOwner,
			DefaultSettings: models.GroupDefaultSettings{
				AllowMemberInvite:   cfg.Group.DefaultSettings.AllowMemberInvite,
				AllowMemberView:     cfg.Group.DefaultSettings.AllowMemberView,
				AllowMemberEditName: cfg.Group.DefaultSettings.AllowMemberEditName,
			},
		}
		service := services.NewGroupService(client, storage, groupConfig)

		// 执行导出操作
		if err := service.ExportGroups(outputFile); err != nil {
			return fmt.Errorf("导出群组数据失败: %w", err)
		}

		fmt.Printf("群组数据已成功导出到: %s\n", outputFile)
		return nil
	},
}

// checkCmd 检查群组命令
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查群组是否存在",
	Long:  "检查指定名称的群组是否已存在",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupName, _ := cmd.Flags().GetString("name")
		if groupName == "" {
			return fmt.Errorf("必须指定群组名称 (--name)")
		}

		// 初始化服务
		client := dingtalk.NewClient(cfg)
		storage := storage.NewFileStorage(cfg.GetDataDir())
		groupConfig := &models.GroupConfig{
			DefaultOwner: cfg.Group.DefaultOwner,
			DefaultSettings: models.GroupDefaultSettings{
				AllowMemberInvite:   cfg.Group.DefaultSettings.AllowMemberInvite,
				AllowMemberView:     cfg.Group.DefaultSettings.AllowMemberView,
				AllowMemberEditName: cfg.Group.DefaultSettings.AllowMemberEditName,
			},
		}
		service := services.NewGroupService(client, storage, groupConfig)

		// 执行检查操作
		exists := service.CheckGroupExists(groupName)
		if exists {
			fmt.Printf("群组 '%s' 已存在\n", groupName)
		} else {
			fmt.Printf("群组 '%s' 不存在\n", groupName)
		}

		return nil
	},
}

func init() {
	// 根命令标志
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件路径")

	// 创建群组命令标志
	createCmd.Flags().StringP("file", "f", "", "CSV文件路径 (必需)")
	createCmd.MarkFlagRequired("file")

	// 添加成员命令标志
	addMemberCmd.Flags().StringP("user-id", "u", "", "用户ID (必需)")
	addMemberCmd.Flags().StringP("group-id", "g", "", "群组ID")
	addMemberCmd.Flags().BoolP("all-groups", "a", false, "添加到所有群组")
	addMemberCmd.MarkFlagRequired("user-id")

	// 移除成员命令标志
	removeMemberCmd.Flags().StringP("user-id", "u", "", "用户ID (必需)")
	removeMemberCmd.Flags().StringP("group-id", "g", "", "群组ID")
	removeMemberCmd.Flags().BoolP("all-groups", "a", false, "从所有群组移除")
	removeMemberCmd.MarkFlagRequired("user-id")

	// 导出命令标志
	exportCmd.Flags().StringP("output", "o", "groups_export.csv", "输出CSV文件路径")

	// 检查命令标志
	checkCmd.Flags().StringP("name", "n", "", "群组名称 (必需)")
	checkCmd.MarkFlagRequired("name")

	// 添加子命令
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addMemberCmd)
	rootCmd.AddCommand(removeMemberCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(checkCmd)
}

func main() {
	// 加载配置
	var err error
	cfg, err = config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令失败: %v\n", err)
		os.Exit(1)
	}
}
