<template>
  <!-- 文件批量处理工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <FolderSearch :size="20" class="text-primary-400" />
        文件批量处理
      </div>
      <div class="page-desc">文件搜索 · 批量复制 · 批量移动 · 批量删除</div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar mb-4">
      <button @click="selectSourceDir" class="btn btn-primary">
        <FolderOpen :size="14" />
        选择目录
      </button>
      <div class="flex-1" />
      <select v-model="fileFilter" class="input-field w-32">
        <option value="*">全部文件</option>
        <option value="*.jpg,*.jpeg,*.png,*.gif,*.bmp,*.webp">图片文件</option>
        <option value="*.pdf">PDF 文档</option>
        <option value="*.doc,*.docx">Word 文档</option>
        <option value="*.txt,*.md">文本文件</option>
        <option value="*.zip,*.rar,*.7z,*.tar,*.gz">压缩文件</option>
        <option value="*.mp4,*.avi,*.mkv,*.mov">视频文件</option>
        <option value="*.mp3,*.wav,*.flac">音频文件</option>
      </select>
      <input v-model="searchPattern" class="input-field w-40" placeholder="搜索文件名..." />
      <button @click="scanFiles" :disabled="!sourceDir" class="btn btn-secondary">
        <Search :size="14" />
        扫描
      </button>
    </div>

    <!-- 主内容区 -->
    <div class="flex-1 flex gap-4 min-h-0">
      <!-- 左侧：文件列表 -->
      <div class="flex-1 flex flex-col min-h-0">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">
            <span>扫描结果</span>
            <span v-if="files.length" class="text-xs opacity-50 ml-2">{{ files.length }} 个文件</span>
          </div>
          <div v-if="files.length" class="flex items-center gap-2">
            <button @click="toggleSelectAll" class="text-xs opacity-60 hover:opacity-100">
              {{ allSelected ? '取消全选' : '全选' }}
            </button>
            <span class="text-xs opacity-50">|</span>
            <span class="text-xs opacity-50">已选 {{ selectedCount }} 个</span>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="!files.length && !scanning" class="flex-1 flex items-center justify-center opacity-30">
          <div class="text-center">
            <FolderSearch :size="40" class="mx-auto mb-2 opacity-50" />
            <div class="text-sm">选择目录并扫描文件</div>
          </div>
        </div>

        <!-- 扫描中 -->
        <div v-if="scanning" class="flex-1 flex items-center justify-center">
          <div class="text-center">
            <Loader2 :size="32" class="mx-auto mb-2 loading-spin text-primary-400" />
            <div class="text-sm opacity-60">正在扫描文件...</div>
          </div>
        </div>

        <!-- 文件列表 -->
        <div v-if="files.length && !scanning" class="flex-1 overflow-auto card p-0">
          <table class="w-full text-sm">
            <thead class="sticky top-0 bg-[var(--bg-card)]">
              <tr class="border-b border-white/5">
                <th class="w-8 p-2"><input type="checkbox" :checked="allSelected" @change="toggleSelectAll" /></th>
                <th class="text-left p-2 opacity-60">文件名</th>
                <th class="text-right p-2 opacity-60 w-20">大小</th>
                <th class="text-right p-2 opacity-60 w-32">修改时间</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(file, idx) in filteredFiles"
                :key="idx"
                class="border-b border-white/5 hover:bg-white/5 cursor-pointer"
                :class="{ 'bg-primary-500/10': file.selected }"
                @click="file.selected = !file.selected"
              >
                <td class="p-2"><input type="checkbox" v-model="file.selected" @click.stop /></td>
                <td class="p-2 truncate">
                  <div class="flex items-center gap-2">
                    <FileIcon :size="14" class="opacity-50 shrink-0" />
                    <span class="truncate">{{ file.name }}</span>
                  </div>
                  <div class="text-xs opacity-40 truncate">{{ file.path }}</div>
                </td>
                <td class="p-2 text-right opacity-60">{{ formatSize(file.size) }}</td>
                <td class="p-2 text-right opacity-60">{{ formatDate(file.modTime) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 右侧：操作面板 -->
      <div class="w-64 flex flex-col gap-3 shrink-0">
        <!-- 操作类型 -->
        <div class="card">
          <div class="font-semibold mb-3 flex items-center gap-2"><Settings2 :size="14" />操作类型</div>
          <div class="space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="operation" value="copy" />
              <Copy :size="14" />
              <span>复制到</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="operation" value="move" />
              <FolderInput :size="14" />
              <span>移动到</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer text-red-400">
              <input type="radio" v-model="operation" value="delete" />
              <Trash2 :size="14" />
              <span>删除</span>
            </label>
          </div>
        </div>

        <!-- 目标目录 -->
        <div v-if="operation !== 'delete'" class="card">
          <div class="font-semibold mb-2">目标目录</div>
          <div class="flex gap-2">
            <input v-model="targetDir" class="input-field flex-1" placeholder="选择目标..." readonly />
            <button @click="selectTargetDir" class="btn btn-secondary px-2">
              <FolderOpen :size="14" />
            </button>
          </div>
        </div>

        <!-- 统计 -->
        <div class="card">
          <div class="font-semibold mb-2">统计信息</div>
          <div class="space-y-1 text-sm">
            <div class="flex justify-between">
              <span class="opacity-60">选中文件</span>
              <span>{{ selectedCount }} 个</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-60">总大小</span>
              <span>{{ formatSize(selectedSize) }}</span>
            </div>
          </div>
        </div>

        <!-- 执行按钮 -->
        <button
          @click="executeOperation"
          :disabled="!canExecute"
          :class="['btn w-full', operation === 'delete' ? 'btn-danger' : 'btn-primary']"
        >
          <Loader2 v-if="isExecuting" :size="14" class="loading-spin" />
          <Play v-else :size="14" />
          {{ isExecuting ? '执行中...' : '执行操作' }}
        </button>

        <!-- 执行日志 -->
        <div v-if="logs.length" class="card flex-1 min-h-0 flex flex-col">
          <div class="flex items-center justify-between mb-2">
            <div class="font-semibold text-sm">执行日志</div>
            <button @click="logs = []" class="text-xs opacity-50">清空</button>
          </div>
          <div class="flex-1 overflow-auto text-xs space-y-1 font-mono">
            <div v-for="(log, idx) in logs" :key="idx" :class="log.type === 'error' ? 'text-red-400' : log.type === 'success' ? 'text-green-400' : 'opacity-60'">
              {{ log.message }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { FolderSearch, FolderOpen, Search, Settings2, Copy, FolderInput, Trash2, Play, Loader2, File as FileIcon } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { OpenFileDialog } from '../../../wailsjs/runtime/runtime'
import { FileOperation } from '../../../wailsjs/go/sysinfo/SysInfo'

const appStore = useAppStore()

// 文件接口
interface FileItem {
  name: string
  path: string
  size: number
  modTime: string
  selected: boolean
}

// 状态
const sourceDir = ref('')
const targetDir = ref('')
const files = ref<FileItem[]>([])
const fileFilter = ref('*')
const searchPattern = ref('')
const operation = ref<'copy' | 'move' | 'delete'>('copy')
const isExecuting = ref(false)
const scanning = ref(false)
const logs = ref<Array<{ type: 'info' | 'success' | 'error'; message: string }>>([])

// 过滤后的文件
const filteredFiles = computed(() => {
  if (!searchPattern.value) return files.value
  const pattern = searchPattern.value.toLowerCase()
  return files.value.filter(f => f.name.toLowerCase().includes(pattern))
})

// 选中数量
const selectedCount = computed(() => files.value.filter(f => f.selected).length)
const selectedSize = computed(() => files.value.filter(f => f.selected).reduce((sum, f) => sum + f.size, 0))
const allSelected = computed(() => files.value.length > 0 && files.value.every(f => f.selected))

// 是否可执行
const canExecute = computed(() => {
  if (isExecuting.value) return false
  if (selectedCount.value === 0) return false
  if (operation.value !== 'delete' && !targetDir.value) return false
  return true
})

// 选择源目录
async function selectSourceDir() {
  try {
    const dir = await OpenFileDialog({ Title: '选择源目录', CanCreateDirectories: true })
    if (dir) {
      sourceDir.value = dir
      appStore.showToast('success', `已选择目录: ${dir}`)
    }
  } catch (e) {
    appStore.showToast('error', '选择目录失败: ' + String(e))
  }
}

// 选择目标目录
async function selectTargetDir() {
  try {
    const dir = await OpenFileDialog({ Title: '选择目标目录', CanCreateDirectories: true })
    if (dir) {
      targetDir.value = dir
      appStore.showToast('success', `已选择目标: ${dir}`)
    }
  } catch (e) {
    appStore.showToast('error', '选择目录失败: ' + String(e))
  }
}

// 扫描文件 - 使用 OpenFileDialog 选择文件
async function scanFiles() {
  if (!sourceDir.value) return
  scanning.value = true
  files.value = []
  try {
    const result = await OpenFileDialog({
      Title: '选择要处理的文件',
      Directory: sourceDir.value,
    })
    if (result) {
      // OpenFileDialog 可能返回单个文件路径或多个文件路径
      const selectedPaths = Array.isArray(result) ? result : [result]
      for (const p of selectedPaths) {
        const name = p.split(/[/\\]/).pop() || p
        files.value.push({
          name,
          path: p,
          size: 0,
          modTime: '',
          selected: false,
        })
      }
      appStore.showToast('success', `已选择 ${files.value.length} 个文件`)
    }
  } catch (e) {
    appStore.showToast('error', '选择文件失败: ' + String(e))
  } finally {
    scanning.value = false
  }
}

// 全选/取消
function toggleSelectAll() {
  const newState = !allSelected.value
  files.value.forEach(f => f.selected = newState)
}

// 执行操作 - 使用 FileOperation 逐个处理
async function executeOperation() {
  if (!canExecute.value) return

  const selectedFiles = files.value.filter(f => f.selected)

  // 删除操作二次确认
  if (operation.value === 'delete') {
    const confirmed = confirm(`确定要删除 ${selectedFiles.length} 个文件吗？此操作不可撤销！`)
    if (!confirmed) return
  }

  isExecuting.value = true
  logs.value = []

  let successCount = 0
  let failCount = 0

  for (const file of selectedFiles) {
    try {
      const dst = operation.value === 'delete' ? '' : (targetDir.value + '/' + file.name)
      await FileOperation(operation.value, file.path, dst)
      logs.value.push({ type: 'success', message: `OK ${file.name}` })
      successCount++
    } catch (err) {
      logs.value.push({ type: 'error', message: `FAIL ${file.name}: ${err}` })
      failCount++
    }
  }

  if (failCount === 0) {
    appStore.showToast('success', `操作完成，成功处理 ${successCount} 个文件`)
    // 如果是删除或移动，从列表移除
    if (operation.value === 'delete' || operation.value === 'move') {
      files.value = files.value.filter(f => !f.selected)
    }
  } else {
    appStore.showToast('warning', `操作完成：${successCount} 成功，${failCount} 失败`)
  }

  isExecuting.value = false
}

// 格式化大小
function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

// 格式化日期
function formatDate(date: string): string {
  if (!date) return '-'
  try {
    const d = new Date(date)
    if (isNaN(d.getTime())) return date
    return d.toLocaleDateString() + ' ' + d.toLocaleTimeString().slice(0, 5)
  } catch {
    return date
  }
}
</script>

<style scoped>
/* 表格样式 */
table {
  border-collapse: collapse;
}

input[type="checkbox"] {
  width: 14px;
  height: 14px;
  cursor: pointer;
}

input[type="radio"] {
  cursor: pointer;
}
</style>
