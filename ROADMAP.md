# XTool 改进计划

## 🔴 必须修复

### 1. 双标题栏问题
- **状态**: ✅ 已修复
- **修改**: `main.go` 中 `DisableFramelessWindowDecorations: false`

### 2. 安装包快捷方式
- **状态**: ✅ 已添加 `installer.nsi`
- **需要**: 安装 NSIS 后运行 `make build-installer`

### 3. 二次安装功能未生效
- **原因**: 前端缓存 + 数据库缓存
- **解决**:
  ```bash
  make clean
  rm -rf frontend/node_modules/.vite
  make build-installer
  ```

## 🟡 建议修复

### 4. DNS 工具国内可用
- **状态**: ✅ 已修复
- **修改**: 添加阿里 DNS、腾讯 DNS、DNSPod、Cloudflare 支持
- **实现**: 支持 DoH wire format 协议，国内默认使用阿里 DNS

### 5. 文件批量处理后端
- **问题**: 当前使用模拟数据
- **方案**: 添加 Go 后端实现真实文件操作

### 6. 应用图标
- **问题**: 缺少自定义图标
- **方案**: 设计并添加 `build/appicon.ico`

## 📋 版本规划

### v1.0.1
- 修复双标题栏
- 修复安装包
- 更新 README

### v1.1.0
- ✅ DNS 工具国内优化 (已完成)
- 文件批量处理后端
- 添加应用图标
- 添加更多主题

### v2.0.0
- 移动端支持（PWA/Capacitor）
- 数据云同步
- 插件系统
