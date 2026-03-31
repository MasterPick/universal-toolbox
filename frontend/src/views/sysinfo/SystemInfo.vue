<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Monitor :size="20" class="text-primary-400"/>系统信息</div>
      <div class="page-desc">查看 CPU、内存、磁盘、网络等系统详情</div>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center opacity-50">
        <div class="loading-spin inline-block mb-2"><RefreshCw :size="24"/></div>
        <div class="text-sm">正在获取系统信息...</div>
      </div>
    </div>

    <div v-else-if="info" class="flex-1 overflow-auto">
      <!-- 基本信息卡 -->
      <div class="grid grid-cols-2 gap-3 mb-4">
        <div class="card">
          <div class="text-xs opacity-50 mb-1">操作系统</div>
          <div class="font-semibold capitalize">{{ info.os }} / {{ info.arch }}</div>
        </div>
        <div class="card">
          <div class="text-xs opacity-50 mb-1">主机名</div>
          <div class="font-semibold">{{ info.hostname }}</div>
        </div>
        <div class="card">
          <div class="text-xs opacity-50 mb-1">运行时间</div>
          <div class="font-semibold">{{ formatUptime(info.uptime) }}</div>
        </div>
        <div class="card">
          <div class="text-xs opacity-50 mb-1">CPU 型号</div>
          <div class="font-semibold text-sm truncate" :title="info.cpuModel">{{ info.cpuModel || 'N/A' }}</div>
        </div>
      </div>

      <!-- 使用率仪表盘 -->
      <div class="grid grid-cols-3 gap-3">
        <div class="card text-center">
          <div class="text-xs opacity-50 mb-2">CPU 使用率</div>
          <div class="text-3xl font-bold" :class="getUsageColor(info.cpuUsage)">
            {{ info.cpuUsage.toFixed(1) }}%
          </div>
          <div class="usage-bar mt-3">
            <div class="usage-fill" :style="{ width: info.cpuUsage + '%' }"
                 :class="getUsageColor(info.cpuUsage)" />
          </div>
          <div class="text-xs opacity-50 mt-1">{{ info.cpuCores }} 核心</div>
        </div>

        <div class="card text-center">
          <div class="text-xs opacity-50 mb-2">内存使用率</div>
          <div class="text-3xl font-bold" :class="getUsageColor(info.memPercent)">
            {{ info.memPercent.toFixed(1) }}%
          </div>
          <div class="usage-bar mt-3">
            <div class="usage-fill" :style="{ width: info.memPercent + '%' }"
                 :class="getUsageColor(info.memPercent)" />
          </div>
          <div class="text-xs opacity-50 mt-1">{{ formatBytes(info.memUsed) }} / {{ formatBytes(info.memTotal) }}</div>
        </div>

        <div class="card text-center">
          <div class="text-xs opacity-50 mb-2">磁盘使用率</div>
          <div class="text-3xl font-bold" :class="getUsageColor(info.diskPercent)">
            {{ info.diskPercent.toFixed(1) }}%
          </div>
          <div class="usage-bar mt-3">
            <div class="usage-fill" :style="{ width: info.diskPercent + '%' }"
                 :class="getUsageColor(info.diskPercent)" />
          </div>
          <div class="text-xs opacity-50 mt-1">{{ formatBytes(info.diskUsed) }} / {{ formatBytes(info.diskTotal) }}</div>
        </div>
      </div>
    </div>

    <!-- 刷新按钮 -->
    <button @click="loadInfo" class="btn btn-secondary mt-4 self-start">
      <RefreshCw :size="14" :class="loading ? 'loading-spin' : ''"/>刷新
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Monitor, RefreshCw } from 'lucide-vue-next'
import { GetSystemInfo } from '../../../wailsjs/go/sysinfo/SysInfo'

interface SysInfoData {
  os: string; arch: string; hostname: string; cpuModel: string
  cpuCores: number; cpuUsage: number; memTotal: number; memUsed: number
  memPercent: number; diskTotal: number; diskUsed: number; diskPercent: number
  uptime: number
}

const info = ref<SysInfoData | null>(null)
const loading = ref(false)

async function loadInfo() {
  loading.value = true
  try {
    info.value = await GetSystemInfo() as SysInfoData
  } finally {
    loading.value = false
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024, sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatUptime(seconds: number): string {
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  return `${d}天 ${h}时 ${m}分`
}

function getUsageColor(pct: number): string {
  if (pct < 60) return 'text-green-400'
  if (pct < 80) return 'text-yellow-400'
  return 'text-red-400'
}

onMounted(loadInfo)
</script>

<style scoped>
.usage-bar {
  height: 6px;
  background: rgba(255,255,255,0.08);
  border-radius: 3px;
  overflow: hidden;
}
.usage-fill {
  height: 100%;
  border-radius: 3px;
  background: currentColor;
  transition: width 0.5s ease;
}
</style>
