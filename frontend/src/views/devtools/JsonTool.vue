<template>
  <!-- JSON 工具页面：格式化、压缩、校验、转义 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Braces :size="20" class="text-primary-400" />
        JSON 工具
      </div>
      <div class="page-desc">格式化 · 压缩 · 校验 · 转义/反转义</div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar mb-3">
      <button @click="formatJSON"    class="btn btn-primary"><Wand2 :size="14"/>格式化</button>
      <button @click="compressJSON"  class="btn btn-secondary"><Minimize2 :size="14"/>压缩</button>
      <button @click="validateJSON"  class="btn btn-secondary"><CheckCircle :size="14"/>校验</button>
      <button @click="escapeJSON"    class="btn btn-secondary"><Code :size="14"/>转义</button>
      <button @click="unescapeJSON"  class="btn btn-secondary"><Undo :size="14"/>反转义</button>
      <div class="flex-1" />
      <button @click="copyOutput"    class="btn btn-secondary"><Copy :size="14"/>复制结果</button>
      <button @click="clearAll"      class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
    </div>

    <!-- 双栏编辑区 -->
    <div class="two-col flex-1 min-h-0">
      <!-- 左侧：输入 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>输入 JSON</span>
          <span class="text-xs opacity-50">{{ inputStats }}</span>
        </div>
        <textarea
          v-model="inputText"
          class="textarea-field flex-1 min-h-0"
          placeholder="在此粘贴 JSON 内容..."
          spellcheck="false"
          @input="onInputChange"
        />
      </div>

      <!-- 右侧：输出 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>处理结果</span>
          <span
            :class="['badge text-xs', resultStatus === 'success' ? 'badge-success' : resultStatus === 'error' ? 'badge-error' : 'badge-info']"
            v-if="resultStatus"
          >
            {{ resultLabel }}
          </span>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto">
          <span v-if="!outputText" class="opacity-30">处理结果将显示在这里...</span>
          <span v-else>{{ outputText }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Braces, Wand2, Minimize2, CheckCircle, Code, Undo, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import {
  FormatJSON, CompressJSON, ValidateJSON, EscapeJSON, UnescapeJSON
} from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const inputText  = ref('')
const outputText = ref('')
const resultStatus = ref<'success' | 'error' | 'info' | ''>('')
const resultLabel  = ref('')

// 实时统计输入字符数
const inputStats = computed(() => {
  if (!inputText.value) return ''
  return `${inputText.value.length} 字符`
})

// 输入变化时自动校验（防抖）
let validateTimer: number
function onInputChange() {
  clearTimeout(validateTimer)
  validateTimer = window.setTimeout(() => {
    if (inputText.value.trim()) autoValidate()
  }, 500)
}

async function autoValidate() {
  try {
    JSON.parse(inputText.value)
    resultStatus.value = 'success'
    resultLabel.value = '✓ 合法 JSON'
  } catch {
    resultStatus.value = 'error'
    resultLabel.value = '✗ 格式错误'
  }
}

async function formatJSON() {
  if (!inputText.value.trim()) return
  const res = await FormatJSON(inputText.value)
  if (res.success) {
    outputText.value = res.data
    resultStatus.value = 'success'
    resultLabel.value = '格式化成功'
    appStore.showToast('success', 'JSON 格式化完成')
  } else {
    outputText.value = res.error
    resultStatus.value = 'error'
    resultLabel.value = '格式化失败'
  }
}

async function compressJSON() {
  if (!inputText.value.trim()) return
  const res = await CompressJSON(inputText.value)
  if (res.success) {
    outputText.value = res.data
    resultStatus.value = 'success'
    resultLabel.value = '压缩成功'
  } else {
    outputText.value = res.error
    resultStatus.value = 'error'
    resultLabel.value = '压缩失败'
  }
}

async function validateJSON() {
  if (!inputText.value.trim()) return
  const res = await ValidateJSON(inputText.value)
  outputText.value = res.data || res.error
  resultStatus.value = res.success ? 'success' : 'error'
  resultLabel.value = res.success ? '✓ 合法' : '✗ 不合法'
  appStore.showToast(res.success ? 'success' : 'error', res.success ? 'JSON 格式正确' : res.error)
}

async function escapeJSON() {
  if (!inputText.value.trim()) return
  const res = await EscapeJSON(inputText.value)
  if (res.success) outputText.value = res.data
}

async function unescapeJSON() {
  if (!inputText.value.trim()) return
  const res = await UnescapeJSON(inputText.value)
  if (res.success) outputText.value = res.data
  else {
    outputText.value = res.error
    resultStatus.value = 'error'
  }
}

async function copyOutput() {
  if (!outputText.value) return
  await navigator.clipboard.writeText(outputText.value)
  appStore.showToast('success', '已复制到剪贴板')
}

function clearAll() {
  inputText.value = ''
  outputText.value = ''
  resultStatus.value = ''
  resultLabel.value = ''
}
</script>
