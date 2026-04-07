<template>
  <div class="page-container">
    <!-- 标题 -->
    <div class="mb-4">
      <div class="page-title">
        <Fingerprint :size="20" class="text-primary-400" />
        UUID 生成器
      </div>
      <div class="page-desc">批量生成唯一标识符 · 多版本支持</div>
    </div>

    <!-- 配置区 -->
    <div class="card mb-4">
      <div class="flex items-center gap-4 flex-wrap">
        <div class="flex items-center gap-2">
          <span class="label mb-0">版本</span>
          <select v-model="version" class="input-field w-20">
            <option value="v4">v4 随机</option>
            <option value="v7">v7 时间</option>
            <option value="v1">v1 时间戳</option>
          </select>
        </div>
        <div class="flex items-center gap-2">
          <span class="label mb-0">数量</span>
          <input v-model.number="count" type="number" min="1" max="1000" class="input-field w-24" />
        </div>
        <div class="flex items-center gap-2">
          <span class="label mb-0">格式</span>
          <select v-model="format" class="input-field w-32">
            <option value="hyphen">带连字符</option>
            <option value="noHyphen">无连字符</option>
            <option value="braces">花括号</option>
            <option value="uppercase">大写</option>
          </select>
        </div>
        <button @click="generate" class="btn btn-primary">
          <RefreshCw :size="14" />
          生成
        </button>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="flex gap-2 mb-4 flex-wrap">
      <button @click="generateSingle" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">
        生成 1 个
      </button>
      <button @click="count = 10; generate()" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">
        生成 10 个
      </button>
      <button @click="count = 50; generate()" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">
        生成 50 个
      </button>
      <button @click="copyAll" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">
        <Copy :size="12" class="inline mr-1" />
        复制全部
      </button>
      <button @click="clear" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">
        <Trash2 :size="12" class="inline mr-1" />
        清空
      </button>
    </div>

    <!-- 结果列表 -->
    <div class="flex-1 flex flex-col min-h-0">
      <!-- 空状态 -->
      <div v-if="!uuids.length" class="flex-1 flex items-center justify-center opacity-30">
        <div class="text-center">
          <Fingerprint :size="32" class="mx-auto mb-2 opacity-50" />
          <div class="text-sm">点击生成按钮创建 UUID</div>
        </div>
      </div>

      <!-- UUID 列表 -->
      <div v-else class="flex-1 overflow-auto space-y-1">
        <div
          v-for="(uuid, idx) in uuids"
          :key="idx"
          class="flex items-center gap-3 p-2 rounded bg-white/5 hover:bg-white/10 group"
        >
          <span class="w-6 h-6 rounded-full bg-primary-500/20 text-primary-400 flex items-center justify-center text-xs">
            {{ idx + 1 }}
          </span>
          <span class="font-mono text-sm flex-1 select-all">{{ uuid }}</span>
          <button
            @click="copyOne(uuid)"
            class="opacity-0 group-hover:opacity-100 transition-opacity text-xs px-2 py-1 rounded bg-white/10"
          >
            <Copy :size="12" />
          </button>
        </div>
      </div>
    </div>

    <!-- 统计 -->
    <div v-if="uuids.length" class="mt-3 text-xs opacity-50">
      共 {{ uuids.length }} 个 UUID · 总字符数 {{ totalChars }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Fingerprint, RefreshCw, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const version = ref('v4')
const count = ref(10)
const format = ref('hyphen')
const uuids = ref<string[]>([])

// 统计
const totalChars = computed(() => uuids.value.join('\n').length)

// 生成 UUID v4
function uuidv4(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

// 生成 UUID v7 (基于时间)
function uuidv7(): string {
  const now = Date.now()
  const hex = now.toString(16).padStart(12, '0')
  const rand = 'xxxx4xxx'.replace(/[x]/g, () => Math.floor(Math.random() * 16).toString(16))
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-7${rand.slice(1, 4)}-${(0x8 | Math.random() * 4).toString(16)}${'x'.repeat(3).replace(/x/g, () => Math.floor(Math.random() * 16).toString(16))}-${'x'.repeat(12).replace(/x/g, () => Math.floor(Math.random() * 16).toString(16))}`
}

// 生成 UUID v1 (模拟时间戳)
function uuidv1(): string {
  const now = Date.now() + 0x01b21dd213814000n // UUID epoch offset
  const hex = now.toString(16).padStart(16, '0')
  const clock = Math.floor(Math.random() * 16384).toString(16).padStart(4, '0')
  const node = Array.from({ length: 12 }, () => Math.floor(Math.random() * 16).toString(16)).join('')
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-1${hex.slice(13, 16)}-${clock}-${node}`
}

// 格式化 UUID
function formatUuid(uuid: string): string {
  switch (format.value) {
    case 'noHyphen':
      return uuid.replace(/-/g, '')
    case 'braces':
      return `{${uuid}}`
    case 'uppercase':
      return uuid.toUpperCase()
    default:
      return uuid
  }
}

// 生成
function generate() {
  const generators: Record<string, () => string> = {
    v4: uuidv4,
    v7: uuidv7,
    v1: uuidv1
  }
  const gen = generators[version.value] || uuidv4
  uuids.value = Array.from({ length: Math.min(count.value, 1000) }, () => formatUuid(gen()))
  appStore.showToast('success', `已生成 ${uuids.value.length} 个 UUID`)
}

// 生成单个
function generateSingle() {
  count.value = 1
  generate()
}

// 复制全部
async function copyAll() {
  if (!uuids.value.length) return
  await navigator.clipboard.writeText(uuids.value.join('\n'))
  appStore.showToast('success', '已复制到剪贴板')
}

// 复制单个
async function copyOne(uuid: string) {
  await navigator.clipboard.writeText(uuid)
  appStore.showToast('success', '已复制')
}

// 清空
function clear() {
  uuids.value = []
}
</script>
