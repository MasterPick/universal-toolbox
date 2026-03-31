// Package log 日志系统模块
// 提供分级日志记录（info/warn/error），支持文件和控制台双输出
package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel 日志级别枚举
type LogLevel int

const (
	INFO  LogLevel = iota // 普通信息
	WARN                  // 警告信息
	ERROR                 // 错误信息
)

// Logger 日志记录器
type Logger struct {
	infoLogger  *log.Logger // info 级别日志记录器
	warnLogger  *log.Logger // warn 级别日志记录器
	errorLogger *log.Logger // error 级别日志记录器
	logFile     *os.File    // 日志文件句柄
}

// NewLogger 创建新的日志记录器实例
// 日志文件存储在用户主目录下的 .xtool/logs/ 文件夹中
func NewLogger() *Logger {
	// 获取用户主目录
	homeDir, _ := os.UserHomeDir()

	// 创建日志目录
	logDir := filepath.Join(homeDir, ".xtool", "logs")
	_ = os.MkdirAll(logDir, 0755)

	// 按日期创建日志文件（每天一个）
	logFileName := fmt.Sprintf("toolbox-%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	// 打开日志文件（追加模式）
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// 若无法写文件，仅输出到控制台
		logFile = os.Stdout
	}

	// 创建各级别日志记录器（同时写文件和控制台）
	flags := log.Ldate | log.Ltime | log.Lshortfile

	return &Logger{
		infoLogger:  log.New(logFile, "[INFO]  ", flags),
		warnLogger:  log.New(logFile, "[WARN]  ", flags),
		errorLogger: log.New(logFile, "[ERROR] ", flags),
		logFile:     logFile,
	}
}

// Info 记录普通信息日志
func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

// Warn 记录警告日志
func (l *Logger) Warn(msg string) {
	l.warnLogger.Println(msg)
}

// Error 记录错误日志
func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.logFile != nil && l.logFile != os.Stdout {
		_ = l.logFile.Close()
	}
}
