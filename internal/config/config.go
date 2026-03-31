// Package config 配置中心模块
// 统一管理应用主题、快捷键、用户偏好等配置项
// 所有配置持久化存储到 SQLite 数据库
package config

import (
	"encoding/json"
	"xtool/internal/db"
)

// AppConfig 应用全局配置结构
type AppConfig struct {
	Theme        string  `json:"theme"`        // 主题：dark/light/auto/blue/green/purple
	Language     string  `json:"language"`     // 语言：zh-CN/en-US
	FontSize     int     `json:"fontSize"`     // 字体大小（px）
	Density      string  `json:"density"`      // 布局密度：compact/normal/spacious
	Transparency float64 `json:"transparency"` // 窗口透明度（0-1）
	AlwaysOnTop  bool    `json:"alwaysOnTop"`  // 窗口置顶
	MinToTray    bool    `json:"minToTray"`    // 最小化到托盘
	AutoBackup   bool    `json:"autoBackup"`   // 自动备份
}

// Config 配置管理器
type Config struct {
	db     *db.Database // 数据库实例
	config *AppConfig   // 当前配置
}

// NewConfig 创建配置管理器并加载已保存的配置
func NewConfig(database *db.Database) *Config {
	c := &Config{
		db: database,
		config: &AppConfig{
			Theme:        "dark",    // 默认深色主题
			Language:     "zh-CN",   // 默认简体中文
			FontSize:     14,        // 默认字体大小
			Density:      "normal",  // 默认正常密度
			Transparency: 0.95,      // 默认透明度
			AlwaysOnTop:  false,     // 默认不置顶
			MinToTray:    true,      // 默认最小化到托盘
			AutoBackup:   true,      // 默认自动备份
		},
	}

	// 从数据库加载已保存的配置
	_ = c.loadFromDB()
	return c
}

// GetConfig 获取当前配置（供前端调用）
func (c *Config) GetConfig() *AppConfig {
	return c.config
}

// SaveConfig 保存配置（供前端调用）
func (c *Config) SaveConfig(configJSON string) error {
	// 解析前端传来的 JSON 配置
	var newConfig AppConfig
	if err := json.Unmarshal([]byte(configJSON), &newConfig); err != nil {
		return err
	}

	// 更新内存中的配置
	c.config = &newConfig

	// 持久化到数据库
	return c.saveToDB()
}

// GetTheme 获取当前主题（供前端调用）
func (c *Config) GetTheme() string {
	return c.config.Theme
}

// SetTheme 设置主题（供前端调用）
func (c *Config) SetTheme(theme string) error {
	c.config.Theme = theme
	return c.saveToDB()
}

// loadFromDB 从数据库加载配置
func (c *Config) loadFromDB() error {
	row := c.db.DB.QueryRow("SELECT value FROM settings WHERE key = 'app_config'")

	var jsonStr string
	if err := row.Scan(&jsonStr); err != nil {
		// 配置不存在，使用默认值
		return nil
	}

	return json.Unmarshal([]byte(jsonStr), c.config)
}

// saveToDB 将配置序列化后存入数据库
func (c *Config) saveToDB() error {
	jsonBytes, err := json.Marshal(c.config)
	if err != nil {
		return err
	}

	// 使用 UPSERT 语法（INSERT OR REPLACE）
	_, err = c.db.DB.Exec(
		"INSERT OR REPLACE INTO settings (key, value, updated_at) VALUES ('app_config', ?, CURRENT_TIMESTAMP)",
		string(jsonBytes),
	)
	return err
}
