// Package db SQLite 数据库初始化与管理模块
// 负责创建数据库连接、初始化表结构、提供数据访问接口
package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite 驱动（CGO）
)

// Database 数据库实例封装
type Database struct {
	DB *sql.DB
}

// Init 初始化 SQLite 数据库
// 数据库文件存储在用户主目录下的 .universal-toolbox/ 文件夹中
func Init() (*Database, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// 创建数据目录（若不存在）
	dataDir := filepath.Join(homeDir, ".universal-toolbox")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// 打开 SQLite 数据库
	dbPath := filepath.Join(dataDir, "toolbox.db")
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// 初始化数据表
	database := &Database{DB: db}
	if err := database.initTables(); err != nil {
		return nil, err
	}

	return database, nil
}

// initTables 创建所有必要的数据表
func (d *Database) initTables() error {
	// 启用外键约束
	if _, err := d.DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	// 代码片段表
	_, err := d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			language TEXT DEFAULT 'text',
			tags TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 备忘录表
	_, err = d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			color TEXT DEFAULT '#ffffff',
			pinned INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 快捷历史记录表（各工具的历史输入）
	_, err = d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tool_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tool_name TEXT NOT NULL,
			input TEXT NOT NULL,
			output TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 用户配置表
	_, err = d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.DB.Close()
}
