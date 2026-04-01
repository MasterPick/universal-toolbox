<template>
  <!-- 主题编辑器页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Palette :size="20" class="text-primary-400" />
        主题编辑器
      </div>
      <div class="page-desc">自定义主题配色 · 实时预览 · 导入导出</div>
    </div>

    <div class="flex-1 flex gap-4 min-h-0">
      <!-- 左侧：颜色配置 -->
      <div class="w-72 flex flex-col gap-3 shrink-0">
        <!-- 基础颜色 -->
        <div class="card">
          <div class="font-semibold mb-3">基础颜色</div>
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <input type="color" v-model="theme.primary" class="w-10 h-8 rounded cursor-pointer" />
              <div class="flex-1">
                <div class="text-sm">主色调</div>
                <input v-model="theme.primary" class="input-field text-xs font-mono" />
              </div>
            </div>
            <div class="flex items-center gap-3">
              <input type="color" v-model="theme.background" class="w-10 h-8 rounded cursor-pointer" />
              <div class="flex-1">
                <div class="text-sm">背景色</div>
                <input v-model="theme.background" class="input-field text-xs font-mono" />
              </div>
            </div>
            <div class="flex items-center gap-3">
              <input type="color" v-model="theme.card" class="w-10 h-8 rounded cursor-pointer" />
              <div class="flex-1">
                <div class="text-sm">卡片背景</div>
                <input v-model="theme.card" class="input-field text-xs font-mono" />
              </div>
            </div>
          </div>
        </div>

        <!-- 文字颜色 -->
        <div class="card">
          <div class="font-semibold mb-3">文字颜色</div>
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <input type="color" v-model="theme.textPrimary" class="w-10 h-8 rounded cursor-pointer" />
              <div class="flex-1">
                <div class="text-sm">主文字</div>
                <input v-model="theme.textPrimary" class="input-field text-xs font-mono" />
              </div>
            </div>
            <div class="flex items-center gap-3">
              <input type="color" v-model="theme.textMuted" class="w-10 h-8 rounded cursor-pointer" />
              <div class="flex-1">
                <div class="text-sm">次要文字</div>
                <input v-model="theme.textMuted" class="input-field text-xs font-mono" />
              </div>
            </div>
          </div>
        </div>

        <!-- 状态颜色 -->
        <div class="card">
          <div class="font-semibold mb-3">状态颜色</div>
          <div class="grid grid-cols-2 gap-2">
            <div class="flex items-center gap-2">
              <input type="color" v-model="theme.success" class="w-8 h-6 rounded cursor-pointer" />
              <span class="text-xs">成功</span>
            </div>
            <div class="flex items-center gap-2">
              <input type="color" v-model="theme.warning" class="w-8 h-6 rounded cursor-pointer" />
              <span class="text-xs">警告</span>
            </div>
            <div class="flex items-center gap-2">
              <input type="color" v-model="theme.error" class="w-8 h-6 rounded cursor-pointer" />
              <span class="text-xs">错误</span>
            </div>
            <div class="flex items-center gap-2">
              <input type="color" v-model="theme.info" class="w-8 h-6 rounded cursor-pointer" />
              <span class="text-xs">信息</span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-2">
          <button @click="resetTheme" class="btn btn-secondary flex-1">
            <RotateCcw :size="14" />
            重置
          </button>
          <button @click="applyTheme" class="btn btn-primary flex-1">
            <Check :size="14" />
            应用
          </button>
        </div>

        <!-- 导入导出 -->
        <div class="card">
          <div class="font-semibold mb-2">导入/导出</div>
          <div class="flex gap-2">
            <button @click="exportTheme" class="btn btn-secondary flex-1 text-xs">
              <Download :size="12" />
              导出
            </button>
            <button @click="importTheme" class="btn btn-secondary flex-1 text-xs">
              <Upload :size="12" />
              导入
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：实时预览 -->
      <div class="flex-1 flex flex-col min-h-0">
        <div class="label mb-2">实时预览</div>
        <div
          class="flex-1 rounded-xl overflow-hidden"
          :style="{ background: theme.background, color: theme.textPrimary }"
        >
          <!-- 预览窗口 -->
          <div class="h-full flex flex-col p-4">
            <!-- 标题栏 -->
            <div class="flex items-center gap-2 mb-4 pb-3 border-b" :style="{ borderColor: theme.primary + '30' }">
              <div class="w-3 h-3 rounded-full" :style="{ background: '#ef4444' }" />
              <div class="w-3 h-3 rounded-full" :style="{ background: '#f59e0b' }" />
              <div class="w-3 h-3 rounded-full" :style="{ background: '#10b981' }" />
              <span class="text-sm ml-2 opacity-60">预览窗口</span>
            </div>

            <!-- 卡片预览 -->
            <div class="grid grid-cols-2 gap-3 mb-4">
              <div class="p-4 rounded-lg" :style="{ background: theme.card }">
                <div class="text-sm font-semibold mb-2">卡片标题</div>
                <div class="text-xs opacity-60 mb-3">这是一段卡片描述文字</div>
                <button class="w-full py-1.5 rounded text-xs text-white" :style="{ background: theme.primary }">
                  主按钮
                </button>
              </div>
              <div class="p-4 rounded-lg" :style="{ background: theme.card }">
                <div class="text-sm font-semibold mb-2">数据统计</div>
                <div class="text-2xl font-bold" :style="{ color: theme.primary }">1,234</div>
                <div class="text-xs opacity-60">较昨日 +12%</div>
              </div>
            </div>

            <!-- 列表预览 -->
            <div class="p-3 rounded-lg flex-1" :style="{ background: theme.card }">
              <div class="text-sm font-semibold mb-2">列表项</div>
              <div class="space-y-2">
                <div v-for="i in 4" :key="i" class="flex items-center gap-3 p-2 rounded" :style="{ background: theme.background }">
                  <div class="w-8 h-8 rounded-full" :style="{ background: theme.primary + '40' }" />
                  <div class="flex-1">
                    <div class="text-sm">列表项 {{ i }}</div>
                    <div class="text-xs" :style="{ color: theme.textMuted }">描述文字</div>
                  </div>
                  <div class="w-12 h-6 rounded text-xs flex items-center justify-center" :style="{ background: theme.primary + '30', color: theme.primary }">
                    标签
                  </div>
                </div>
              </div>
            </div>

            <!-- 状态标签 -->
            <div class="flex gap-2 mt-4">
              <span class="px-2 py-1 rounded text-xs" :style="{ background: theme.success + '20', color: theme.success }">成功</span>
              <span class="px-2 py-1 rounded text-xs" :style="{ background: theme.warning + '20', color: theme.warning }">警告</span>
              <span class="px-2 py-1 rounded text-xs" :style="{ background: theme.error + '20', color: theme.error }">错误</span>
              <span class="px-2 py-1 rounded text-xs" :style="{ background: theme.info + '20', color: theme.info }">信息</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Palette, RotateCcw, Check, Download, Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 主题配置
const theme = ref({
  primary: '#6366f1',
  background: '#0f0f17',
  card: '#1a1a2e',
  textPrimary: '#e2e8f0',
  textMuted: '#94a3b8',
  success: '#10b981',
  warning: '#f59e0b',
  error: '#ef4444',
  info: '#3b82f6',
})

// 默认主题（用于重置）
const defaultTheme = { ...theme.value }

// 重置主题
function resetTheme() {
  theme.value = { ...defaultTheme }
  appStore.showToast('success', '已重置为默认主题')
}

// 应用主题
function applyTheme() {
  // 生成 CSS 变量
  const css = `
    :root {
      --accent: ${theme.value.primary};
      --bg-primary: ${theme.value.background};
      --bg-card: ${theme.value.card};
      --text-primary: ${theme.value.textPrimary};
      --text-muted: ${theme.value.textMuted};
      --success: ${theme.value.success};
      --warning: ${theme.value.warning};
      --danger: ${theme.value.error};
    }
  `

  // 注入样式
  let styleEl = document.getElementById('custom-theme')
  if (!styleEl) {
    styleEl = document.createElement('style')
    styleEl.id = 'custom-theme'
    document.head.appendChild(styleEl)
  }
  styleEl.textContent = css

  appStore.showToast('success', '主题已应用')
}

// 导出主题
function exportTheme() {
  const json = JSON.stringify(theme.value, null, 2)
  const blob = new Blob([json], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'theme.json'
  a.click()
  URL.revokeObjectURL(url)
  appStore.showToast('success', '主题已导出')
}

// 导入主题
function importTheme() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    try {
      const text = await file.text()
      const data = JSON.parse(text)
      theme.value = { ...defaultTheme, ...data }
      appStore.showToast('success', '主题已导入')
    } catch {
      appStore.showToast('error', '导入失败：文件格式错误')
    }
  }
  input.click()
}
</script>

<style scoped>
input[type="color"] {
  border: none;
  padding: 0;
  background: transparent;
}

input[type="color"]::-webkit-color-swatch-wrapper {
  padding: 0;
}

input[type="color"]::-webkit-color-swatch {
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
