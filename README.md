<div align="center">

<img src="docs/logo.png" alt="XTool Logo" width="120" height="120">

# XTool

**一站式全平台桌面工具集，为开发者与高级用户赋能**

[![Version](https://img.shields.io/badge/version-2.0.0-blue?style=flat-square)](https://github.com/MasterPick/universal-toolbox/releases/latest)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-brightgreen?style=flat-square)]()
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go)]()
[![Wails](https://img.shields.io/badge/Wails-v2-61dafb?style=flat-square)]()
[![Vue](https://img.shields.io/badge/Vue-3.4-42b883?style=flat-square&logo=vue.js)]()

[立即下载](#-快速开始) · [功能总览](#-功能总览) · [技术架构](#-技术架构) · [贡献指南](#-贡献指南)

</div>

---

## 核心亮点

| 特性 | 说明 |
|:----:|------|
| **极速启动** | Go 原生编译，启动时间 < 500ms，即开即用 |
| **纯本地安全** | 零网络请求，零数据上传，所有计算在本地完成 |
| **跨平台** | Windows / macOS (Intel + Apple Silicon) / Linux 全平台覆盖 |
| **功能全面** | 40+ 实用工具，覆盖开发、系统、日常、网络四大场景 |
| **轻量便携** | 单文件 < 25MB，无需安装，解压即用 |
| **界面现代** | Windows 11 云母效果，深色/浅色/多主题，响应式布局 |

---

## 功能总览

### 开发工具

| 工具 | 功能描述 | 标签 |
|------|----------|------|
| **JSON 工具** | 格式化、压缩、校验、转义/反转义 | `格式化` `校验` |
| **JSON 对比** | 两个 JSON 差异对比，高亮显示 | `对比` `差异` |
| **XML 工具** | XML 格式化美化 | `格式化` |
| **YAML 工具** | YAML 格式化与验证 | `格式化` `验证` |
| **Base64** | 文本/文件 Base64 编码解码 | `编码` `解码` |
| **URL 编解码** | URL 编码与解码 | `编码` `解码` |
| **哈希计算** | MD5 / SHA1 / SHA256 / SHA512 | `加密` `哈希` |
| **加密解密** | AES / DES 加密解密 | `加密` `安全` |
| **文本工具** | 字符统计、大小写转换、去重 | `文本` `统计` |
| **二维码** | 生成/解析二维码，支持 Logo | `生成` `解析` |
| **UUID 生成** | 批量生成 UUID v1/v4 | `生成` `批量` |
| **时间戳** | 时间戳与日期时间互转 | `转换` `时间` |
| **正则测试** | 正则表达式实时匹配测试 | `匹配` `测试` |
| **代码片段** | 保存管理常用代码片段 | `管理` `收藏` |

### 系统工具

| 工具 | 功能描述 | 标签 |
|------|----------|------|
| **进程管理** | 查看进程列表，强制终止进程 | `监控` `管理` |
| **端口查看** | 查看端口占用，一键释放端口 | `网络` `诊断` |
| **系统信息** | CPU/内存/磁盘实时监控 | `监控` `硬件` |
| **批量重命名** | 文件批量重命名，支持正则表达式 | `文件` `批量` |
| **取色器** | 屏幕取色，8 种颜色格式转换 | `颜色` `设计` |
| **图片工具** | 格式转换(PNG/JPG/WebP)、压缩、调整尺寸 | `图片` `转换` |
| **文件批量处理** | 批量复制/移动/删除文件 | `文件` `批量` |

### 日常工具

| 工具 | 功能描述 | 标签 |
|------|----------|------|
| **计算器** | 标准 + 科学计算器 | `数学` `科学` |
| **单位换算** | 长度/重量/温度/速度/时间换算 | `换算` `单位` |
| **备忘录** | 轻量级便签管理 | `笔记` `管理` |

### 网络工具

| 工具 | 功能描述 | 标签 |
|------|----------|------|
| **Ping 测试** | 主机连通性检测，延迟统计 | `诊断` `连通性` |
| **内网扫描** | 局域网设备发现 | `扫描` `发现` |
| **HTTP 测试** | HTTP 接口请求测试 | `API` `调试` |
| **DNS 查询** | DNS 记录查询(A/AAAA/MX/TXT/NS/CNAME) | `DNS` `查询` |
| **Hosts 编辑** | 系统 Hosts 文件可视化管理 | `编辑` `管理` |

### 设置

| 功能 | 描述 |
|------|------|
| **主题编辑器** | 自定义主题配色，实时预览 |
| **快捷键管理** | 自定义快捷键绑定 |

---

## v2.0.0 新增功能

本次版本新增了 10 个工具，进一步扩展了 XTool 的能力边界：

| 新增工具 | 分类 | 说明 |
|----------|------|------|
| **JSON 对比** | 开发工具 | 两个 JSON 差异对比，高亮显示不同之处 |
| **YAML 工具** | 开发工具 | YAML 格式化与语法验证 |
| **加密解密** | 开发工具 | AES / DES 对称加密解密 |
| **取色器** | 系统工具 | 屏幕取色，支持 8 种颜色格式输出 |
| **图片工具** | 系统工具 | PNG/JPG/WebP 格式转换、压缩、尺寸调整 |
| **文件批量处理** | 系统工具 | 批量复制、移动、删除文件 |
| **DNS 查询** | 网络工具 | DNS 记录查询（A/AAAA/MX/TXT/NS/CNAME） |
| **主题编辑器** | 设置 | 自定义主题配色方案，实时预览效果 |
| **快捷键管理** | 设置 | 自定义各工具快捷键绑定 |
| **macOS arm64** | 平台 | 新增 Apple Silicon (M1/M2/M3) 原生支持 |

---

## 快速开始

### 环境要求

| 工具 | 版本要求 |
|------|----------|
| Go | 1.24+ |
| Node.js | 22+ |
| Wails CLI | v2.12+ |

### 安装与开发

```bash
# 1. 克隆仓库
git clone https://github.com/MasterPick/universal-toolbox.git
cd universal-toolbox

# 2. 安装构建依赖
make install-deps

# 3. 启动开发模式（热重载）
make dev

# 4. 构建当前平台
make all

# 5. 构建指定平台安装包
make build-windows      # Windows .exe
make build-linux        # Linux .deb + tar.gz
make build-macos        # macOS Intel .dmg
make build-macos-arm64  # macOS Apple Silicon .dmg

# 6. 统一打包
make package
```

### 代码检查

```bash
# 运行 Go vet + 前端类型检查
make lint
```

---

## 技术架构

```
XTool
|
+-- 后端 (Go)
|   +-- Wails v2          跨平台桌面应用框架
|   +-- gopsutil          系统信息采集
|   +-- SQLite            本地数据持久化
|   +-- CGO               SQLite 编译支持
|   +-- 模块化设计        各功能独立、可插拔
|
+-- 前端 (Vue 3)
|   +-- Vue 3.4           响应式框架
|   +-- TypeScript        类型安全
|   +-- TailwindCSS       原子化 CSS
|   +-- Pinia             状态管理
|   +-- Vue Router        路由管理
|   +-- Lucide Icons      图标库
|
+-- 构建与发布
    +-- Makefile          本地多平台构建
    +-- GoReleaser        自动化发布
    +-- GitHub Actions    CI/CD 多平台并行构建
```

---

## 界面预览

<div align="center">

| 深色主题 | 浅色主题 |
|:--------:|:--------:|
| ![Dark Theme](docs/screenshot-dark.png) | ![Light Theme](docs/screenshot-light.png) |

</div>

---

## 支持平台

| 平台 | 架构 | 安装包格式 | 最低版本 |
|------|------|-----------|----------|
| **Windows** | amd64 | `.exe` 单文件 + NSIS 安装包 | Windows 10 |
| **Linux** | amd64 | `.deb` + `.tar.gz` | Ubuntu 20.04+ / Debian 10+ |
| **macOS** | amd64 (Intel) | `.dmg` | macOS 10.15 Catalina+ |
| **macOS** | arm64 (Apple Silicon) | `.dmg` | macOS 11 Big Sur+ |

> Windows 11 推荐使用，可体验原生云母效果。

---

## 贡献指南

欢迎参与 XTool 的开发与改进！

1. **Fork** 本仓库
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'feat: add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 创建 **Pull Request**

提交规范建议遵循 [Conventional Commits](https://www.conventionalcommits.org/)：
- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `refactor:` 代码重构
- `perf:` 性能优化

---

## 许可证

本项目基于 [MIT License](LICENSE) 开源。

---

<div align="center">

**如果这个项目对你有帮助，请给一个 Star**

Made with passion by [MasterPick](https://github.com/MasterPick)

</div>
