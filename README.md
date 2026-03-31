# XTool

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-brightgreen)
![License](https://img.shields.io/badge/license-MIT-green)
![Go](https://img.shields.io/badge/Go-1.22-00ADD8)
![Vue](https://img.shields.io/badge/Vue-3.4-42b883)

**一站式桌面工具集，功能丰富、启动极快、界面现代、动效流畅、无广告**

[下载最新版本](#下载) · [功能介绍](#功能) · [开发文档](#开发)

</div>

---

## ✨ 特性

- 🚀 **极速启动**：基于 Go + Wails，启动速度 < 800ms
- 💎 **现代界面**：支持云母/亚克力效果，深色/浅色/多主题切换
- 🛠️ **功能丰富**：30+ 实用工具，覆盖开发、系统、日常、网络
- 🔒 **纯本地**：不收集用户数据，不联网，不留后门
- 📦 **体积小**：打包体积 < 30MB，内存占用 < 120MB
- 🌍 **跨平台**：支持 Windows / Linux / macOS

## 📦 下载

| 平台 | 下载 |
|------|------|
| Windows 10/11 | [XTool.exe](https://github.com/MasterPick/xtool/releases/latest) |
| Linux (amd64) | [xtool_linux_amd64.tar.gz](https://github.com/MasterPick/xtool/releases/latest) |
| macOS | [xtool_macos.dmg](https://github.com/MasterPick/xtool/releases/latest) |

## 🛠️ 功能清单

### 开发工具
| 功能 | 描述 |
|------|------|
| JSON 工具 | 格式化、压缩、校验、转义/反转义 |
| XML 工具 | XML 格式化美化 |
| Base64 | 文本 Base64 编码/解码 |
| URL 编解码 | URL 编码与解码 |
| 哈希计算 | MD5 / SHA1 / SHA256 |
| 文本工具 | 对比、查找替换、字符统计 |
| 二维码 | 生成二维码 |
| UUID 生成 | 批量生成 UUID v4 |
| 时间戳 | 时间戳 ↔ 日期时间互转 |
| 正则测试 | 正则表达式实时测试 |
| 代码片段 | 保存管理常用代码片段 |

### 系统工具
| 功能 | 描述 |
|------|------|
| 进程管理 | 查看/终止进程 |
| 端口查看 | 查看端口占用，一键释放 |
| 系统信息 | CPU/内存/磁盘使用率 |
| 批量重命名 | 文件批量重命名（支持正则） |

### 日常工具
| 功能 | 描述 |
|------|------|
| 计算器 | 标准 + 科学计算 |
| 单位换算 | 长度/重量/温度/速度 |
| 备忘录 | 便签式备忘管理 |

### 网络工具
| 功能 | 描述 |
|------|------|
| Ping 测试 | 主机连通性检测 |
| 内网扫描 | 局域网设备发现 |
| HTTP 测试 | 接口请求与响应查看 |
| Hosts 编辑 | 系统 Hosts 文件编辑 |

## 🏗️ 技术架构

```
XTool
├── 后端 (Go)
│   ├── Wails v2         # 跨平台桌面框架
│   ├── gopsutil         # 系统信息获取
│   ├── SQLite           # 本地数据存储
│   └── 各功能模块       # 独立、可插拔
└── 前端 (Vue3)
    ├── Vue 3.4          # 响应式框架
    ├── TypeScript       # 类型安全
    ├── TailwindCSS      # 原子化 CSS
    ├── Pinia            # 状态管理
    └── Vue Router       # 路由管理
```

## 💻 本地开发

### 环境要求
- Go 1.22+
- Node.js 22+
- Wails CLI v2.12+

### 安装步骤

```bash
# 1. 克隆仓库
git clone https://github.com/MasterPick/xtool.git
cd xtool

# 2. 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 3. 安装前端依赖
cd frontend && npm install && cd ..

# 4. 开发模式运行（热重载）
wails dev

# 5. 构建生产版本
wails build
```

## 📝 开发规范

- 所有 Go 代码包含完整中文注释
- Vue 组件遵循 Composition API 规范
- 错误全部捕获并友好提示
- 日志分级：info / warn / error

## 📄 许可证

[MIT License](LICENSE) © 2026 MasterPick
