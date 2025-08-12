package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	DingTalk DingTalkConfig `mapstructure:"dingtalk"`
	App      AppConfig      `mapstructure:"app"`
	Group    GroupConfig    `mapstructure:"group"`
}

// DingTalkConfig 钉钉应用配置
type DingTalkConfig struct {
	AppKey      string `mapstructure:"app_key"`
	AppSecret   string `mapstructure:"app_secret"`
	AccessToken string `mapstructure:"access_token"`
	CorpID      string `mapstructure:"corp_id"`
	BaseURL     string `mapstructure:"base_url"`
}

// AppConfig 应用配置
type AppConfig struct {
	DataDir  string `mapstructure:"data_dir"`
	LogLevel string `mapstructure:"log_level"`
	Debug    bool   `mapstructure:"debug"`
}

// GroupConfig 群组配置
type GroupConfig struct {
	DefaultOwner    string               `mapstructure:"default_owner"`
	DefaultSettings GroupDefaultSettings `mapstructure:"default_settings"`
}

// GroupDefaultSettings 群组默认设置
type GroupDefaultSettings struct {
	AllowMemberInvite   bool `mapstructure:"allow_member_invite"`
	AllowMemberView     bool `mapstructure:"allow_member_view"`
	AllowMemberEditName bool `mapstructure:"allow_member_edit_name"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 如果指定了配置文件路径，使用指定路径
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 否则在默认位置查找
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.ti-dding")
		viper.AddConfigPath("/etc/ti-dding")
	}

	// 设置环境变量前缀
	viper.SetEnvPrefix("TI_DDING")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认配置
			fmt.Println("配置文件未找到，使用默认配置")
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	viper.SetDefault("dingtalk.base_url", "https://oapi.dingtalk.com")
	viper.SetDefault("app.data_dir", "./data")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("group.default_settings.allow_member_invite", true)
	viper.SetDefault("group.default_settings.allow_member_view", true)
	viper.SetDefault("group.default_settings.allow_member_edit_name", false)
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	// 验证钉钉配置
	if config.DingTalk.AppKey == "" && config.DingTalk.AccessToken == "" {
		return fmt.Errorf("钉钉配置不完整：需要提供 AppKey/AppSecret 或 AccessToken")
	}

	// 验证数据目录
	if config.App.DataDir == "" {
		return fmt.Errorf("数据目录不能为空")
	}

	// 确保数据目录存在
	if err := ensureDataDir(config.App.DataDir); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	return nil
}

// ensureDataDir 确保数据目录存在
func ensureDataDir(dataDir string) error {
	absPath, err := filepath.Abs(dataDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return err
	}

	return nil
}

// GetAccessToken 获取访问令牌
func (c *Config) GetAccessToken() string {
	if c.DingTalk.AccessToken != "" {
		return c.DingTalk.AccessToken
	}
	// 如果没有直接提供访问令牌，需要通过 AppKey/AppSecret 获取
	// 这里可以添加获取访问令牌的逻辑
	return ""
}

// IsDebug 是否启用调试模式
func (c *Config) IsDebug() bool {
	return c.App.Debug
}

// GetDataDir 获取数据目录
func (c *Config) GetDataDir() string {
	return c.App.DataDir
}

// GetLogLevel 获取日志级别
func (c *Config) GetLogLevel() string {
	return c.App.LogLevel
}
