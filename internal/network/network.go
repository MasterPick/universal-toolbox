// Package network 网络工具模块
// 提供 Ping 测试、IP 扫描、HTTP 接口测试、Hosts 编辑等网络工具
package network

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// NetworkTools 网络工具结构体（Wails 绑定到前端）
type NetworkTools struct {
	httpClient *http.Client // HTTP 客户端（可复用）
}

// PingResult Ping 测试结果
type PingResult struct {
	Host    string  `json:"host"`    // 目标主机
	Alive   bool    `json:"alive"`   // 是否在线
	LatencyMs float64 `json:"latencyMs"` // 延迟（毫秒）
	Error   string  `json:"error"`   // 错误信息
}

// ScanResult 内网扫描结果
type ScanResult struct {
	IP      string `json:"ip"`      // IP 地址
	Online  bool   `json:"online"`  // 是否在线
	Latency int64  `json:"latency"` // 延迟（毫秒）
}

// HTTPTestResult HTTP 接口测试结果
type HTTPTestResult struct {
	StatusCode  int               `json:"statusCode"`  // HTTP 状态码
	Status      string            `json:"status"`      // 状态描述
	LatencyMs   float64           `json:"latencyMs"`   // 请求耗时（毫秒）
	Body        string            `json:"body"`        // 响应体
	Headers     map[string]string `json:"headers"`     // 响应头
	ContentType string            `json:"contentType"` // 内容类型
	Error       string            `json:"error"`       // 错误信息
}

// NewNetworkTools 创建网络工具模块实例
func NewNetworkTools() *NetworkTools {
	return &NetworkTools{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // 30 秒超时
		},
	}
}

// ============================================================
// Ping 工具
// ============================================================

// PingHost 对指定主机进行 Ping 测试（通过 TCP 连接检测，无需 root 权限）
func (n *NetworkTools) PingHost(host string) PingResult {
	result := PingResult{Host: host}

	// 使用 TCP 连接到常见端口来模拟 Ping（不需要 ICMP 权限）
	start := time.Now()
	conn, err := net.DialTimeout("tcp", host+":80", 3*time.Second)
	if err != nil {
		// 尝试其他端口
		conn, err = net.DialTimeout("tcp", host+":443", 3*time.Second)
		if err != nil {
			// 最后尝试 DNS 解析
			_, dnsErr := net.LookupHost(host)
			if dnsErr != nil {
				result.Alive = false
				result.Error = "主机不可达：" + err.Error()
				return result
			}
			// DNS 可以解析但 TCP 不通
			result.Alive = true // 主机存在但端口未开放
		}
	}

	if conn != nil {
		conn.Close()
	}

	result.Alive = true
	result.LatencyMs = float64(time.Since(start).Nanoseconds()) / 1e6 // 转毫秒

	return result
}

// PingMultiple 批量 Ping 多个主机
func (n *NetworkTools) PingMultiple(hosts []string) []PingResult {
	results := make([]PingResult, len(hosts))
	var wg sync.WaitGroup

	for i, host := range hosts {
		wg.Add(1)
		go func(idx int, h string) {
			defer wg.Done()
			results[idx] = n.PingHost(h)
		}(i, host)
	}

	wg.Wait()
	return results
}

// ============================================================
// 内网扫描工具
// ============================================================

// ScanLAN 扫描内网 IP 段，返回在线设备列表
// subnet: 如 "192.168.1"，将扫描 192.168.1.1-254
func (n *NetworkTools) ScanLAN(subnet string) []ScanResult {
	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 并发扫描，限制并发数为 50
	semaphore := make(chan struct{}, 50)

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", subnet, i)
		wg.Add(1)
		semaphore <- struct{}{}

		go func(targetIP string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			start := time.Now()
			conn, err := net.DialTimeout("tcp", targetIP+":80", 500*time.Millisecond)
			if err != nil {
				conn, err = net.DialTimeout("tcp", targetIP+":22", 500*time.Millisecond)
			}

			if err == nil {
				conn.Close()
				mu.Lock()
				results = append(results, ScanResult{
					IP:      targetIP,
					Online:  true,
					Latency: time.Since(start).Milliseconds(),
				})
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	return results
}

// GetLocalSubnet 获取本机所在子网（用于扫描建议）
func (n *NetworkTools) GetLocalSubnet() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				// 返回前三段（如 192.168.1）
				parts := strings.Split(ip4.String(), ".")
				if len(parts) >= 3 {
					return strings.Join(parts[:3], "."), nil
				}
			}
		}
	}

	return "192.168.1", nil // 默认值
}

// ============================================================
// HTTP 接口测试工具
// ============================================================

// HTTPRequest 发送 HTTP 请求并返回结果
// method: GET/POST/PUT/DELETE/PATCH
func (n *NetworkTools) HTTPRequest(method, url, body string, headers map[string]string) HTTPTestResult {
	result := HTTPTestResult{}

	// 创建请求
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		result.Error = "请求创建失败：" + err.Error()
		return result
	}

	// 设置默认请求头
	req.Header.Set("User-Agent", "XTool/1.0")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置自定义请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发送请求并计时
	start := time.Now()
	resp, err := n.httpClient.Do(req)
	if err != nil {
		result.Error = "请求失败：" + err.Error()
		return result
	}
	defer resp.Body.Close()

	result.LatencyMs = float64(time.Since(start).Nanoseconds()) / 1e6
	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.ContentType = resp.Header.Get("Content-Type")

	// 读取响应体（限制最多 1MB）
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		result.Error = "读取响应失败：" + err.Error()
		return result
	}
	result.Body = string(respBody)

	// 收集响应头
	result.Headers = make(map[string]string)
	for k, v := range resp.Header {
		result.Headers[k] = strings.Join(v, ", ")
	}

	return result
}

// ============================================================
// Hosts 文件编辑工具
// ============================================================

// GetHostsContent 读取 hosts 文件内容
func (n *NetworkTools) GetHostsContent() (string, error) {
	hostsPath := n.getHostsPath()
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		return "", fmt.Errorf("读取 hosts 文件失败：%v", err)
	}
	return string(content), nil
}

// SaveHostsContent 保存 hosts 文件内容（需要管理员权限）
func (n *NetworkTools) SaveHostsContent(content string) error {
	hostsPath := n.getHostsPath()
	// 创建备份
	backupPath := hostsPath + ".bak"
	original, _ := os.ReadFile(hostsPath)
	_ = os.WriteFile(backupPath, original, 0644)

	// 写入新内容
	if err := os.WriteFile(hostsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("保存 hosts 文件失败（可能需要管理员权限）：%v", err)
	}
	return nil
}

// getHostsPath 获取当前系统的 hosts 文件路径
func (n *NetworkTools) getHostsPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}
