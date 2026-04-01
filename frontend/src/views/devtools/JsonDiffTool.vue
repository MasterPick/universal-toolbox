<template>
  <!-- JSON 对比工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <GitCompare :size="20" class="text-primary-400" />
        JSON 对比
      </div>
      <div class="page-desc">对比两个 JSON 的差异 · 高亮显示变化</div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar mb-4">
      <button @click="swapJson" class="btn btn-secondary">
        <ArrowLeftRight :size="14" />
        交换
      </button>
      <button @click="clearAll" class="btn btn-secondary">
        <Trash2 :size="14" />
        清空
      </button>
      <div class="flex-1" />
      <select v-model="diffMode" class="input-field w-28">
        <option value="line">按行对比</option>
        <option value="key">按键对比</option>
      </select>
      <button @click="compare" :disabled="!canCompare" class="btn btn-primary">
        <GitCompare :size="14" />
        开始对比
      </button>
    </div>

    <!-- 输入区域 -->
    <div class="two-col flex-1 min-h-0 mb-4">
      <!-- 左侧 JSON -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span class="text-blue-400">JSON A（原始）</span>
          <span class="text-xs opacity-50">{{ jsonAStats }}</span>
        </div>
        <textarea
          v-model="jsonA"
          class="textarea-field flex-1 min-h-0 font-mono text-sm"
          placeholder="粘贴原始 JSON..."
          spellcheck="false"
        />
      </div>

      <!-- 右侧 JSON -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span class="text-green-400">JSON B（修改后）</span>
          <span class="text-xs opacity-50">{{ jsonBStats }}</span>
        </div>
        <textarea
          v-model="jsonB"
          class="textarea-field flex-1 min-h-0 font-mono text-sm"
          placeholder="粘贴修改后的 JSON..."
          spellcheck="false"
        />
      </div>
    </div>

    <!-- 对比结果 -->
    <div v-if="diffResult.length" class="card flex-1 min-h-0 flex flex-col">
      <div class="flex items-center justify-between mb-3">
        <div class="font-semibold flex items-center gap-2">
          <FileDiff :size="16" />
          对比结果
        </div>
        <div class="flex items-center gap-4 text-xs">
          <span class="flex items-center gap-1">
            <span class="w-3 h-3 rounded bg-green-500/30" />
            新增
          </span>
          <span class="flex items-center gap-1">
            <span class="w-3 h-3 rounded bg-red-500/30" />
            删除
          </span>
          <span class="flex items-center gap-1">
            <span class="w-3 h-3 rounded bg-yellow-500/30" />
            修改
          </span>
        </div>
      </div>

      <!-- 统计 -->
      <div class="flex gap-4 mb-3 text-sm">
        <span class="text-green-400">+{{ stats.added }} 新增</span>
        <span class="text-red-400">-{{ stats.removed }} 删除</span>
        <span class="text-yellow-400">~{{ stats.modified }} 修改</span>
        <span class="opacity-50">= 相同</span>
      </div>

      <!-- 差异列表 -->
      <div class="flex-1 overflow-auto font-mono text-sm">
        <div
          v-for="(diff, idx) in diffResult"
          :key="idx"
          :class="[
            'diff-line px-3 py-1 border-l-2',
            diff.type === 'add' ? 'bg-green-500/10 border-green-500' :
            diff.type === 'remove' ? 'bg-red-500/10 border-red-500' :
            diff.type === 'modify' ? 'bg-yellow-500/10 border-yellow-500' :
            'border-transparent opacity-50'
          ]"
        >
          <span v-if="diff.type === 'add'" class="text-green-400 mr-2">+</span>
          <span v-else-if="diff.type === 'remove'" class="text-red-400 mr-2">-</span>
          <span v-else-if="diff.type === 'modify'" class="text-yellow-400 mr-2">~</span>
          <span v-else class="opacity-30 mr-2"> </span>
          <span class="opacity-70">{{ diff.path }}</span>
          <span v-if="diff.type === 'modify'" class="text-xs ml-2">
            <span class="text-red-400 line-through">{{ diff.oldValue }}</span>
            <span class="mx-1">→</span>
            <span class="text-green-400">{{ diff.newValue }}</span>
          </span>
          <span v-else-if="diff.type === 'add'" class="text-xs ml-2 text-green-400">{{ diff.value }}</span>
          <span v-else-if="diff.type === 'remove'" class="text-xs ml-2 text-red-400 line-through">{{ diff.value }}</span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="jsonA || jsonB" class="flex-1 flex items-center justify-center opacity-30">
      <div class="text-center">
        <GitCompare :size="32" class="mx-auto mb-2 opacity-50" />
        <div class="text-sm">点击"开始对比"查看差异</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { GitCompare, ArrowLeftRight, Trash2, FileDiff } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const jsonA = ref('')
const jsonB = ref('')
const diffMode = ref<'line' | 'key'>('key')
const diffResult = ref<Array<{ type: 'add' | 'remove' | 'modify' | 'same'; path: string; value?: any; oldValue?: any; newValue?: any }>>([])

// 统计
const stats = ref({ added: 0, removed: 0, modified: 0 })

// 输入统计
const jsonAStats = computed(() => jsonA.value ? `${jsonA.value.length} 字符` : '')
const jsonBStats = computed(() => jsonB.value ? `${jsonB.value.length} 字符` : '')

// 是否可以对比
const canCompare = computed(() => jsonA.value.trim() && jsonB.value.trim())

// 交换
function swapJson() {
  const temp = jsonA.value
  jsonA.value = jsonB.value
  jsonB.value = temp
}

// 清空
function clearAll() {
  jsonA.value = ''
  jsonB.value = ''
  diffResult.value = []
  stats.value = { added: 0, removed: 0, modified: 0 }
}

// 对比
function compare() {
  diffResult.value = []
  stats.value = { added: 0, removed: 0, modified: 0 }

  try {
    const objA = JSON.parse(jsonA.value)
    const objB = JSON.parse(jsonB.value)

    if (diffMode.value === 'key') {
      compareByKey(objA, objB, '')
    } else {
      compareByLine()
    }

    appStore.showToast('success', `对比完成：${stats.value.added} 新增, ${stats.value.removed} 删除, ${stats.value.modified} 修改`)
  } catch (err) {
    appStore.showToast('error', 'JSON 格式错误: ' + String(err))
  }
}

// 按键对比
function compareByKey(a: any, b: any, path: string) {
  const allKeys = new Set([...Object.keys(a || {}), ...Object.keys(b || {})])

  for (const key of allKeys) {
    const currentPath = path ? `${path}.${key}` : key
    const hasA = a && key in a
    const hasB = b && key in b

    if (!hasA && hasB) {
      diffResult.value.push({ type: 'add', path: currentPath, value: JSON.stringify(b[key]) })
      stats.value.added++
    } else if (hasA && !hasB) {
      diffResult.value.push({ type: 'remove', path: currentPath, value: JSON.stringify(a[key]) })
      stats.value.removed++
    } else if (typeof a[key] === 'object' && typeof b[key] === 'object' && a[key] !== null && b[key] !== null) {
      compareByKey(a[key], b[key], currentPath)
    } else if (JSON.stringify(a[key]) !== JSON.stringify(b[key])) {
      diffResult.value.push({ type: 'modify', path: currentPath, oldValue: JSON.stringify(a[key]), newValue: JSON.stringify(b[key]) })
      stats.value.modified++
    } else {
      diffResult.value.push({ type: 'same', path: currentPath, value: JSON.stringify(a[key]) })
    }
  }
}

// 按行对比
function compareByLine() {
  const linesA = jsonA.value.split('\n')
  const linesB = jsonB.value.split('\n')
  const maxLen = Math.max(linesA.length, linesB.length)

  for (let i = 0; i < maxLen; i++) {
    const lineA = linesA[i] || ''
    const lineB = linesB[i] || ''

    if (!linesA[i] && linesB[i]) {
      diffResult.value.push({ type: 'add', path: `行 ${i + 1}`, value: lineB })
      stats.value.added++
    } else if (linesA[i] && !linesB[i]) {
      diffResult.value.push({ type: 'remove', path: `行 ${i + 1}`, value: lineA })
      stats.value.removed++
    } else if (lineA !== lineB) {
      diffResult.value.push({ type: 'modify', path: `行 ${i + 1}`, oldValue: lineA, newValue: lineB })
      stats.value.modified++
    } else {
      diffResult.value.push({ type: 'same', path: `行 ${i + 1}`, value: lineA })
    }
  }
}
</script>

<style scoped>
.diff-line {
  line-height: 1.5;
}

.diff-line:hover {
  background: rgba(255, 255, 255, 0.03);
}
</style>
