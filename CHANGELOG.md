# Changelog

所有版本的重要变更记录。

格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)

## [2.0.0] - 2026-04-10

### 新增
- **JSON 对比工具**：两个 JSON 差异对比，高亮显示不同之处
- **YAML 工具**：YAML 格式化与语法验证
- **加密解密工具**：AES / DES 对称加密解密
- **取色器**：屏幕取色，支持 8 种颜色格式转换（HEX/RGB/HSL/CMYK 等）
- **图片工具**：PNG/JPG/WebP 格式转换、压缩、尺寸调整
- **文件批量处理**：批量复制、移动、删除文件
- **DNS 查询工具**：DNS 记录查询（A/AAAA/MX/TXT/NS/CNAME）
- **主题编辑器**：自定义主题配色方案，实时预览效果
- **快捷键管理**：自定义各工具快捷键绑定
- **macOS arm64 原生支持**：新增 Apple Silicon (M1/M2/M3) 构建与 DMG 打包

### 优化
- 修复 JSON 格式化工具在处理大文件时的性能问题
- 修复 Base64 编解码工具对二进制文件的处理异常
- 修复端口查看工具在部分 Linux 发行版上的权限问题
- 修复计算器科学模式下的精度丢失问题
- 修复备忘录内容在极端情况下丢失的问题
- 修复 HTTP 测试工具对 HTTPS 证书校验的处理
- 修复进程管理工具在 macOS 上的显示异常
- 修复正则测试工具在特定正则表达式下的崩溃问题
- 优化系统信息采集的兼容性，支持更多硬件平台
- 优化窗口启动速度，减少白屏时间

### 打包
- Makefile 新增 `build-macos-arm64` 目标，支持 Apple Silicon 原生构建
- Makefile 新增 `package` 统一打包目标和 `lint` 代码检查目标
- Makefile 添加版本号注入（-ldflags）和 CGO_ENABLED=1
- .goreleaser.yaml 修复 binary 命名（`万能工具箱` -> `xtool`），统一项目命名
- .goreleaser.yaml 添加 NSIS 安装包配置（使用项目自定义图标）
- .goreleaser.yaml 添加 .deb 包 depends 字段（libgtk-3-0, libwebkit2gtk-4.0-37）
- .goreleaser.yaml 添加 macOS universal binary 配置
- GitHub Actions NSIS 使用项目自定义图标（build/windows/icon.ico）
- GitHub Actions 修复仓库 URL 不一致问题（统一为 MasterPick/universal-toolbox）
- GitHub Actions 添加版本号自动回写（发布时自动更新 wails.json）
- wails.json 版本号更新为 2.0.0

### 文档
- README.md 完全重写，采用专业高端风格
- README.md 新增核心亮点、功能标签、v2.0.0 新增功能、支持平台等章节
- README.md 更新技术架构图，添加构建与发布模块
- CHANGELOG.md 添加 v2.0.0 版本记录

## [1.1.0] - 2026-03-31

### 新增
- 🚀 **多平台安装包支持**
  - Windows：`.exe` 单文件 + NSIS 安装向导（含桌面快捷方式、控制面板卸载）
  - Linux：`.deb` 安装包（含应用菜单图标）+ `tar.gz` 通用压缩包
  - macOS：Intel `.dmg` + Apple Silicon (M1/M2/M3) `.dmg` 双架构
- ⚡ **GitHub Actions CI/CD 升级**
  - 三平台并行构建，互不依赖
  - 自动生成 SHA256 校验和文件
  - Release 页面展示详细安装说明
- 📦 **Makefile 本地构建支持**
  - `make dev` 启动热重载开发模式
  - `make build-linux` 生成 .deb + tar.gz
  - `make build-macos` 生成双架构 .dmg
- 📋 **平台构建文档**：各平台详细构建说明

### 优化
- GitHub Release 页面美化，含功能展开列表
- 构建产物 SHA256 校验和汇总到单文件

## [1.0.0] - 2026-03-31

### 新增
- 🛠️ **开发工具模块**：JSON 格式化/压缩/校验/转义，XML 格式化，Base64 编解码，URL 编解码，哈希计算（MD5/SHA1/SHA256），文本对比/替换/统计，UUID 批量生成，时间戳转换，正则测试，代码片段管理
- 💻 **系统工具模块**：进程管理（查看/终止），端口占用查看/释放，系统信息（CPU/内存/磁盘），批量文件重命名
- 📱 **日常工具模块**：标准+科学计算器，长度/重量/温度/速度换算，备忘录
- 🌐 **网络工具模块**：Ping 测试，内网 IP 扫描，HTTP 接口测试，Hosts 文件编辑
- 🎨 **主题系统**：深色/浅色/自动/蓝/绿/紫/橙/云母效果 8 种主题
- ⚙️ **设置中心**：主题切换、字体大小、窗口行为配置
- 🪟 **Windows 11 云母效果**：原生半透明背景

### 技术栈
- Go 1.21 + Wails v2.9
- Vue 3.4 + TypeScript + TailwindCSS
- SQLite（本地数据）
- gopsutil（系统信息）
