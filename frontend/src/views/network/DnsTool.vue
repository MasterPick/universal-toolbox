<template>
  <!-- DNS 查询工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Globe2 :size="20" class="text-primary-400" />
        DNS 查询
      </div>
      <div class="page-desc">域名解析 · DNS 记录查询 · 批量查询</div>
    </div>

    <!-- 查询输入 -->
    <div class="card mb-4">
      <div class="flex gap-3">
        <input
          v-model="domain"
          class="input-field flex-1"
          placeholder="输入域名，如 example.com"
          @keyup.enter="query"
        />
        <select v-model="recordType" class="input-field w-28">
          <option value="A">A 记录</option>
          <option value="AAAA">AAAA 记录</option>
          <option value="CNAME">CNAME</option>
          <option value="MX">MX 记录</option>
          <option value="TXT">TXT 记录</option>
          <option value="NS">NS 记录</option>
          <option value="SOA">SOA</option>
        </select>
        <button @click="query" :disabled="isLoading" class="btn btn-primary">
          <Loader2 v-if="isLoading" :size="14" class="loading-spin" />
          <Search v-else :size="14" />
          查询
        </button>
      </div>
    </div>

    <!-- 快捷域名 -->
    <div class="flex gap-2 mb-4 flex-wrap">
      <span class="text-xs opacity-50">快捷：</span>
      <button
        v-for="d in quickDomains"
        :key="d"
        @click="domain = d; query()"
        class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10"
      >
        {{ d }}
      </button>
    </div>

    <!-- 查询结果 -->
    <div class="flex-1 flex flex-col min-h-0">
      <!-- 空状态 -->
      <div v-if="!result && !isLoading" class="flex-1 flex items-center justify-center opacity-30">
        <div class="text-center">
          <Globe2 :size="40" class="mx-auto mb-2 opacity-50" />
          <div class="text-sm">输入域名进行 DNS 查询</div>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-else-if="isLoading" class="flex-1 flex items-center justify-center">
        <div class="text-center">
          <Loader2 :size="32" class="mx-auto mb-2 loading-spin text-primary-400" />
          <div class="text-sm opacity-60">正在查询...</div>
        </div>
      </div>

      <!-- 结果展示 -->
      <div v-else-if="result" class="flex-1 overflow-auto">
        <!-- 错误 -->
        <div v-if="result.error" class="card bg-red-500/10 border-red-500/30">
          <div class="flex items-center gap-2 text-red-400">
            <AlertCircle :size="16" />
            <span>{{ result.error }}</span>
          </div>
        </div>

        <!-- 成功结果 -->
        <div v-else class="space-y-3">
          <!-- 基本信息 -->
          <div class="card">
            <div class="flex items-center gap-2 mb-3">
              <CheckCircle :size="16" class="text-green-400" />
              <span class="font-semibold">{{ result.domain }}</span>
              <span class="badge badge-info">{{ result.type }} 记录</span>
            </div>
            <div class="text-xs opacity-50">
              查询时间: {{ result.queryTime }}ms · DNS 服务器: {{ result.dnsServer || '系统默认' }}
            </div>
          </div>

          <!-- 记录列表 -->
          <div class="card">
            <div class="font-semibold mb-3">解析结果 ({{ result.records?.length || 0 }} 条)</div>
            <div class="space-y-2">
              <div
                v-for="(record, idx) in result.records"
                :key="idx"
                class="flex items-center gap-3 p-3 rounded-lg bg-white/5"
              >
                <div class="w-8 h-8 rounded-full flex items-center justify-center" :class="getRecordColor(result.type)">
                  {{ idx + 1 }}
                </div>
                <div class="flex-1 font-mono text-sm">
                  {{ record.value || record }}
                </div>
                <div v-if="record.ttl" class="text-xs opacity-50">
                  TTL: {{ record.ttl }}
                </div>
                <button @click="copyValue(record.value || record)" class="btn btn-secondary text-xs py-1">
                  <Copy :size="12" />
                </button>
              </div>
            </div>
          </div>

          <!-- 原始响应 -->
          <div class="card">
            <div class="flex items-center justify-between mb-2">
              <div class="font-semibold text-sm">原始响应</div>
              <button @click="showRaw = !showRaw" class="text-xs opacity-50">
                {{ showRaw ? '收起' : '展开' }}
              </button>
            </div>
            <div v-if="showRaw" class="code-output text-xs">
              <pre>{{ JSON.stringify(result.raw, null, 2) }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 查询历史 -->
    <div v-if="history.length" class="mt-4">
      <div class="flex items-center justify-between mb-2">
        <div class="text-sm opacity-60">查询历史</div>
        <button @click="history = []" class="text-xs opacity-50">清空</button>
      </div>
      <div class="flex gap-2 flex-wrap">
        <button
          v-for="h in history.slice(0, 10)"
          :key="h"
          @click="domain = h; query()"
          class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10 font-mono"
        >
          {{ h }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Globe2, Search, Loader2, AlertCircle, CheckCircle, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const domain = ref('')
const recordType = ref('A')
const isLoading = ref(false)
const showRaw = ref(false)
const result = ref<any>(null)
const history = ref<string[]>([])

// 快捷域名
const quickDomains = ['google.com', 'github.com', 'baidu.com', 'qq.com', 'openai.com']

// 查询 DNS
async function query() {
  if (!domain.value.trim()) return

  isLoading.value = true
  result.value = null

  try {
    // 使用 DNS over HTTPS (DoH) 服务
    const dohUrl = `https://dns.google/resolve?name=${domain.value}&type=${recordType.value}`
    const response = await fetch(dohUrl)
    const data = await response.json()

    if (data.Status !== 0) {
      result.value = {
        error: getDnsStatus(data.Status),
        domain: domain.value,
        type: recordType.value
      }
    } else {
      const records = data.Answer?.map((a: any) => ({
        value: a.data,
        ttl: a.TTL
      })) || []

      result.value = {
        domain: domain.value,
        type: recordType.value,
        records,
        queryTime: 0,
        dnsServer: 'dns.google (8.8.8.8)',
        raw: data
      }

      // 添加到历史
      if (!history.value.includes(domain.value)) {
        history.value.unshift(domain.value)
        if (history.value.length > 20) history.value.pop()
        localStorage.setItem('dnsHistory', JSON.stringify(history.value))
      }
    }
  } catch (err) {
    result.value = {
      error: '查询失败: ' + String(err),
      domain: domain.value,
      type: recordType.value
    }
  } finally {
    isLoading.value = false
  }
}

// DNS 状态码转换
function getDnsStatus(status: number): string {
  const statusMap: Record<number, string> = {
    0: '成功',
    1: '格式错误',
    2: '服务器错误',
    3: '域名不存在',
    4: '未实现',
    5: '拒绝'
  }
  return statusMap[status] || `未知错误 (${status})`
}

// 记录类型颜色
function getRecordColor(type: string): string {
  const colors: Record<string, string> = {
    A: 'bg-blue-500/20 text-blue-400',
    AAAA: 'bg-green-500/20 text-green-400',
    CNAME: 'bg-purple-500/20 text-purple-400',
    MX: 'bg-yellow-500/20 text-yellow-400',
    TXT: 'bg-orange-500/20 text-orange-400',
    NS: 'bg-cyan-500/20 text-cyan-400',
    SOA: 'bg-pink-500/20 text-pink-400'
  }
  return colors[type] || 'bg-gray-500/20 text-gray-400'
}

// 复制值
async function copyValue(value: string) {
  await navigator.clipboard.writeText(value)
  appStore.showToast('success', '已复制')
}

// 加载历史
onMounted(() => {
  const saved = localStorage.getItem('dnsHistory')
  if (saved) {
    try {
      history.value = JSON.parse(saved)
    } catch {
      // 忽略
    }
  }
})
</script>
