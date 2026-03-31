// Package devtools 开发工具模块
// 提供 JSON/XML/YAML 格式化、Base64、URL 编解码、哈希计算、UUID 生成等开发常用工具
package devtools

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
	"xtool/internal/db"

	"github.com/google/uuid"
)

// DevTools 开发工具结构体（Wails 绑定到前端）
type DevTools struct {
	db *db.Database // 数据库连接，用于存储代码片段和历史记录
}

// Snippet 代码片段结构
type Snippet struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ToolResult 工具执行结果通用结构
type ToolResult struct {
	Success bool   `json:"success"` // 是否成功
	Data    string `json:"data"`    // 结果数据
	Error   string `json:"error"`   // 错误信息（失败时填充）
}

// NewDevTools 创建开发工具模块实例
func NewDevTools(database *db.Database) *DevTools {
	return &DevTools{db: database}
}

// ============================================================
// JSON 工具
// ============================================================

// FormatJSON 格式化 JSON 字符串（美化输出，带缩进）
func (d *DevTools) FormatJSON(input string) ToolResult {
	// 解析 JSON 验证合法性
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ToolResult{Success: false, Error: "JSON 格式错误：" + err.Error()}
	}

	// 重新序列化（带 2 空格缩进）
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ToolResult{Success: false, Error: "格式化失败：" + err.Error()}
	}

	return ToolResult{Success: true, Data: string(formatted)}
}

// CompressJSON 压缩 JSON 字符串（去除多余空格和换行）
func (d *DevTools) CompressJSON(input string) ToolResult {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ToolResult{Success: false, Error: "JSON 格式错误：" + err.Error()}
	}

	compressed, err := json.Marshal(obj)
	if err != nil {
		return ToolResult{Success: false, Error: "压缩失败：" + err.Error()}
	}

	return ToolResult{Success: true, Data: string(compressed)}
}

// ValidateJSON 校验 JSON 是否合法
func (d *DevTools) ValidateJSON(input string) ToolResult {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ToolResult{Success: false, Error: "不合法的 JSON：" + err.Error()}
	}
	return ToolResult{Success: true, Data: "✓ JSON 格式正确"}
}

// EscapeJSON 转义 JSON 字符串（用于嵌入到字符串字面量中）
func (d *DevTools) EscapeJSON(input string) ToolResult {
	escaped, err := json.Marshal(input)
	if err != nil {
		return ToolResult{Success: false, Error: "转义失败：" + err.Error()}
	}
	// 去掉首尾引号
	result := string(escaped)
	result = result[1 : len(result)-1]
	return ToolResult{Success: true, Data: result}
}

// UnescapeJSON 反转义 JSON 字符串
func (d *DevTools) UnescapeJSON(input string) ToolResult {
	// 加上引号以构成合法的 JSON 字符串值
	quoted := `"` + input + `"`
	var result string
	if err := json.Unmarshal([]byte(quoted), &result); err != nil {
		return ToolResult{Success: false, Error: "反转义失败：" + err.Error()}
	}
	return ToolResult{Success: true, Data: result}
}

// ============================================================
// XML 工具
// ============================================================

// FormatXML 格式化 XML 字符串
func (d *DevTools) FormatXML(input string) ToolResult {
	// 解码后重新编码（带缩进）
	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = false

	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		if err := encoder.EncodeToken(token); err != nil {
			return ToolResult{Success: false, Error: "XML 编码失败：" + err.Error()}
		}
	}

	if err := encoder.Flush(); err != nil {
		return ToolResult{Success: false, Error: "格式化失败：" + err.Error()}
	}

	result := buf.String()
	if result == "" {
		return ToolResult{Success: false, Error: "XML 格式错误，无法解析"}
	}

	return ToolResult{Success: true, Data: result}
}

// ============================================================
// Base64 工具
// ============================================================

// Base64Encode 将字符串编码为 Base64
func (d *DevTools) Base64Encode(input string) ToolResult {
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return ToolResult{Success: true, Data: encoded}
}

// Base64Decode 将 Base64 字符串解码
func (d *DevTools) Base64Decode(input string) ToolResult {
	// 去除可能的空格和换行
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")

	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		// 尝试 URL 安全的 Base64
		decoded, err = base64.URLEncoding.DecodeString(input)
		if err != nil {
			return ToolResult{Success: false, Error: "Base64 解码失败：" + err.Error()}
		}
	}

	// 检查是否为合法 UTF-8 文本
	if !utf8.Valid(decoded) {
		// 以十六进制显示
		return ToolResult{Success: true, Data: fmt.Sprintf("(二进制数据) 十六进制: %x", decoded)}
	}

	return ToolResult{Success: true, Data: string(decoded)}
}

// ============================================================
// URL 编解码工具
// ============================================================

// URLEncode URL 编码字符串
func (d *DevTools) URLEncode(input string) ToolResult {
	encoded := url.QueryEscape(input)
	return ToolResult{Success: true, Data: encoded}
}

// URLDecode URL 解码字符串
func (d *DevTools) URLDecode(input string) ToolResult {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return ToolResult{Success: false, Error: "URL 解码失败：" + err.Error()}
	}
	return ToolResult{Success: true, Data: decoded}
}

// ============================================================
// 哈希计算工具
// ============================================================

// CalcMD5 计算字符串的 MD5 哈希值
func (d *DevTools) CalcMD5(input string) ToolResult {
	hash := md5.Sum([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// CalcSHA1 计算字符串的 SHA1 哈希值
func (d *DevTools) CalcSHA1(input string) ToolResult {
	hash := sha1.Sum([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// CalcSHA256 计算字符串的 SHA256 哈希值
func (d *DevTools) CalcSHA256(input string) ToolResult {
	hash := sha256.Sum256([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// ============================================================
// 文本工具
// ============================================================

// TextCompare 对比两个文本，返回差异行
func (d *DevTools) TextCompare(text1, text2 string) ToolResult {
	lines1 := strings.Split(text1, "\n")
	lines2 := strings.Split(text2, "\n")

	var sb strings.Builder
	maxLen := len(lines1)
	if len(lines2) > maxLen {
		maxLen = len(lines2)
	}

	diffCount := 0
	for i := 0; i < maxLen; i++ {
		var l1, l2 string
		if i < len(lines1) {
			l1 = lines1[i]
		}
		if i < len(lines2) {
			l2 = lines2[i]
		}
		if l1 != l2 {
			diffCount++
			sb.WriteString(fmt.Sprintf("第 %d 行:\n  原文: %s\n  新文: %s\n\n", i+1, l1, l2))
		}
	}

	if diffCount == 0 {
		return ToolResult{Success: true, Data: "✓ 两段文本完全相同"}
	}

	return ToolResult{
		Success: true,
		Data:    fmt.Sprintf("共发现 %d 处差异：\n\n%s", diffCount, sb.String()),
	}
}

// TextReplace 在文本中执行查找替换
func (d *DevTools) TextReplace(text, search, replace string, useRegex bool) ToolResult {
	if text == "" || search == "" {
		return ToolResult{Success: false, Error: "文本和搜索内容不能为空"}
	}

	var result string
	if useRegex {
		// 使用正则表达式替换
		re, err := regexp.Compile(search)
		if err != nil {
			return ToolResult{Success: false, Error: "正则表达式错误：" + err.Error()}
		}
		result = re.ReplaceAllString(text, replace)
	} else {
		// 普通字符串替换
		result = strings.ReplaceAll(text, search, replace)
	}

	return ToolResult{Success: true, Data: result}
}

// TextStats 统计文本的字符数、单词数、行数等信息
func (d *DevTools) TextStats(text string) ToolResult {
	charCount := len([]rune(text))       // 字符数（支持中文）
	byteCount := len(text)               // 字节数
	lineCount := strings.Count(text, "\n") + 1 // 行数
	wordCount := len(strings.Fields(text))     // 英文单词数

	// 统计中文字符数
	chineseCount := 0
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			chineseCount++
		}
	}

	result := fmt.Sprintf(
		"字符数（含空格）: %d\n字节数: %d\n行数: %d\n英文单词数: %d\n中文字符数: %d",
		charCount, byteCount, lineCount, wordCount, chineseCount,
	)

	return ToolResult{Success: true, Data: result}
}

// ============================================================
// UUID 生成工具
// ============================================================

// GenerateUUID 生成一个 UUID v4
func (d *DevTools) GenerateUUID() ToolResult {
	return ToolResult{Success: true, Data: uuid.New().String()}
}

// GenerateUUIDs 批量生成多个 UUID
func (d *DevTools) GenerateUUIDs(count int) ToolResult {
	if count <= 0 || count > 100 {
		count = 10 // 默认 10 个，最多 100 个
	}

	var ids []string
	for i := 0; i < count; i++ {
		ids = append(ids, uuid.New().String())
	}

	return ToolResult{Success: true, Data: strings.Join(ids, "\n")}
}

// ============================================================
// 时间戳工具
// ============================================================

// TimestampToDatetime 时间戳（秒/毫秒）转日期时间字符串
func (d *DevTools) TimestampToDatetime(timestamp int64) ToolResult {
	var t time.Time

	// 自动判断是秒还是毫秒级时间戳
	if timestamp > 1e12 {
		// 毫秒级时间戳
		t = time.Unix(timestamp/1000, (timestamp%1000)*int64(time.Millisecond)).Local()
	} else {
		// 秒级时间戳
		t = time.Unix(timestamp, 0).Local()
	}

	result := fmt.Sprintf(
		"本地时间: %s\nUTC 时间: %s\n时区: %s",
		t.Format("2006-01-02 15:04:05"),
		t.UTC().Format("2006-01-02 15:04:05 UTC"),
		t.Format("MST"),
	)

	return ToolResult{Success: true, Data: result}
}

// DatetimeToTimestamp 日期时间字符串转时间戳
func (d *DevTools) DatetimeToTimestamp(datetime string) ToolResult {
	// 支持多种格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.ParseInLocation(format, datetime, time.Local)
		if err == nil {
			break
		}
	}

	if err != nil {
		return ToolResult{Success: false, Error: "日期格式不正确，请使用 YYYY-MM-DD HH:MM:SS 格式"}
	}

	result := fmt.Sprintf(
		"秒级时间戳: %d\n毫秒级时间戳: %d",
		t.Unix(),
		t.UnixMilli(),
	)

	return ToolResult{Success: true, Data: result}
}

// GetCurrentTimestamp 获取当前时间戳
func (d *DevTools) GetCurrentTimestamp() ToolResult {
	now := time.Now()
	result := fmt.Sprintf(
		"当前时间: %s\n秒级时间戳: %d\n毫秒级时间戳: %d",
		now.Format("2006-01-02 15:04:05"),
		now.Unix(),
		now.UnixMilli(),
	)
	return ToolResult{Success: true, Data: result}
}

// ============================================================
// 正则表达式工具
// ============================================================

// RegexTest 测试正则表达式匹配结果
func (d *DevTools) RegexTest(pattern, text string) ToolResult {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ToolResult{Success: false, Error: "正则表达式语法错误：" + err.Error()}
	}

	matches := re.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return ToolResult{Success: true, Data: "未找到匹配项"}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("找到 %d 个匹配项：\n\n", len(matches)))

	for i, match := range matches {
		sb.WriteString(fmt.Sprintf("匹配 #%d: %s\n", i+1, match[0]))
		for j, group := range match[1:] {
			sb.WriteString(fmt.Sprintf("  捕获组 %d: %s\n", j+1, group))
		}
	}

	return ToolResult{Success: true, Data: sb.String()}
}

// ============================================================
// 代码片段管理
// ============================================================

// GetSnippets 获取所有代码片段
func (d *DevTools) GetSnippets() ([]Snippet, error) {
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, language, tags, created_at, updated_at FROM snippets ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Language, &s.Tags, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}

// SaveSnippet 保存代码片段（新增或更新）
func (d *DevTools) SaveSnippet(title, content, language, tags string) (int64, error) {
	result, err := d.db.DB.Exec(
		"INSERT INTO snippets (title, content, language, tags) VALUES (?, ?, ?, ?)",
		title, content, language, tags,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// DeleteSnippet 删除指定代码片段
func (d *DevTools) DeleteSnippet(id int64) error {
	_, err := d.db.DB.Exec("DELETE FROM snippets WHERE id = ?", id)
	return err
}

// SearchSnippets 搜索代码片段（标题、内容、标签）
func (d *DevTools) SearchSnippets(keyword string) ([]Snippet, error) {
	query := "%" + keyword + "%"
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, language, tags, created_at, updated_at FROM snippets WHERE title LIKE ? OR content LIKE ? OR tags LIKE ? ORDER BY updated_at DESC",
		query, query, query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Language, &s.Tags, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}
