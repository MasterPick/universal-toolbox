<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Clock :size="20" class="text-primary-400"/>时间戳转换</div>
      <div class="page-desc">时间戳与日期时间互相转换</div>
    </div>

    <!-- 当前时间 -->
    <div class="card mb-4">
      <div class="flex items-center justify-between">
        <span class="text-sm opacity-60">当前时间</span>
        <button @click="refreshNow" class="btn btn-secondary py-1 px-2 text-xs"><RefreshCw :size="12"/>刷新</button>
      </div>
      <div class="code-output mt-2 text-sm">{{ nowInfo }}</div>
    </div>

    <div class="two-col">
      <!-- 时间戳 → 日期 -->
      <div class="card">
        <div class="label mb-3">时间戳 → 日期时间</div>
        <input v-model="tsInput" class="input-field mb-3"
          placeholder="输入时间戳（秒/毫秒）..." />
        <button @click="tsToDate" class="btn btn-primary w-full mb-3">转换</button>
        <div class="code-output text-sm">{{ tsResult || '结果显示在这里...' }}</div>
      </div>

      <!-- 日期 → 时间戳 -->
      <div class="card">
        <div class="label mb-3">日期时间 → 时间戳</div>
        <input v-model="dateInput" class="input-field mb-3"
          placeholder="如：2024-01-01 12:00:00" />
        <button @click="dateToTs" class="btn btn-primary w-full mb-3">转换</button>
        <div class="code-output text-sm">{{ dateResult || '结果显示在这里...' }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Clock, RefreshCw } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { TimestampToDatetime, DatetimeToTimestamp, GetCurrentTimestamp } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const tsInput   = ref('')
const dateInput = ref('')
const tsResult  = ref('')
const dateResult = ref('')
const nowInfo   = ref('')

async function refreshNow() {
  const res = await GetCurrentTimestamp()
  nowInfo.value = res.data
}

async function tsToDate() {
  if (!tsInput.value) return
  const ts = parseInt(tsInput.value)
  if (isNaN(ts)) { appStore.showToast('error', '请输入有效的时间戳'); return }
  const res = await TimestampToDatetime(ts)
  tsResult.value = res.success ? res.data : res.error
}

async function dateToTs() {
  if (!dateInput.value) return
  const res = await DatetimeToTimestamp(dateInput.value)
  if (res.success) {
    dateResult.value = res.data
  } else {
    appStore.showToast('error', res.error)
    dateResult.value = res.error
  }
}

onMounted(refreshNow)
</script>
