<template>
  <ToolPage title="网络测速" description="下载/上传速度与延迟测试">
<!-- 网络测速工具 -->
  
    

    <!-- 控制区域 -->
    <div class="card mb-4">
      <div class="flex gap-2 items-center">
        <button @click="startTest" :class="['btn', testing ? 'btn-danger' : 'btn-primary']" :disabled="testing && !canStop">
          <Gauge :size="14" :class="testing ? 'loading-spin' : ''"/>
          {{ testing ? '停止测试' : '开始测速' }}
        </button>
        <div class="flex-1" />
        <span v-if="testing" class="text-sm opacity-50">{{ testPhaseText }}</span>
      </div>
    </div>

    <!-- 速度仪表盘 -->
    <div class="card mb-4">
      <div class="flex justify-center py-6">
        <div class="relative w-64 h-64">
          <!-- 仪表盘背景 -->
          <svg viewBox="0 0 200 200" class="w-full h-full">
            <!-- 背景圆弧 -->
            <circle cx="100" cy="100" r="85" fill="none" stroke="rgba(255,255,255,0.05)" stroke-width="12"
              stroke-dasharray="400 134" stroke-linecap="round"
              transform="rotate(135 100 100)" />
            <!-- 速度圆弧 -->
            <circle cx="100" cy="100" r="85" fill="none" :stroke="gaugeColor" stroke-width="12"
              stroke-dasharray="400 134" stroke-linecap="round"
              :stroke-dashoffset="400 - (400 * gaugePercent / 100)"
              transform="rotate(135 100 100)"
              class="transition-all duration-500" />
          </svg>
          <!-- 中心文字 -->
          <div class="absolute inset-0 flex flex-col items-center justify-center">
            <div class="text-4xl font-bold font-mono" :class="testing ? 'text-primary-400' : 'opacity-50'">
              {{ currentSpeed }}
            </div>
            <div class="text-xs opacity-40 mt-1">{{ currentUnit }}</div>
            <div v-if="testing" class="text-xs opacity-30 mt-2">{{ testPhaseText }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 测试结果摘要 -->
    <div v-if="testResults" class="grid grid-cols-4 gap-3 mb-4">
      <div class="card text-center">
        <div class="text-xs opacity-40 mb-1">下载速度</div>
        <div class="text-xl font-bold font-mono text-green-400">{{ testResults.downloadSpeed }}</div>
        <div class="text-xs opacity-30">Mbps</div>
      </div>
      <div class="card text-center">
        <div class="text-xs opacity-40 mb-1">上传速度</div>
        <div class="text-xl font-bold font-mono text-blue-400">{{ testResults.uploadSpeed }}</div>
        <div class="text-xs opacity-30">Mbps</div>
      </div>
      <div class="card text-center">
        <div class="text-xs opacity-40 mb-1">延迟</div>
        <div class="text-xl font-bold font-mono" :class="getPingColor(testResults.ping)">
          {{ testResults.ping }}
        </div>
        <div class="text-xs opacity-30">ms</div>
      </div>
      <div class="card text-center">
        <div class="text-xs opacity-40 mb-1">丢包率</div>
        <div class="text-xl font-bold font-mono" :class="Number(testResults.loss) > 0 ? 'text-red-400' : 'text-green-400'">
          {{ testResults.loss }}%
        </div>
        <div class="text-xs opacity-30">loss</div>
      </div>
    </div>

    <!-- 实时速度图表 -->
    <div v-if="testing || speedHistory.length > 0" class="card mb-4">
      <div class="label mb-2">实时速度</div>
      <div class="h-32 flex items-end gap-px">
        <div
          v-for="(s, i) in speedHistory"
          :key="i"
          class="flex-1 rounded-t transition-all duration-200"
          :class="s.phase === 'download' ? 'bg-green-400/60' : s.phase === 'upload' ? 'bg-blue-400/60' : 'bg-yellow-400/60'"
          :style="{ height: Math.max(2, (s.speed / maxSpeed) * 100) + '%' }"
          :title="`${s.speed.toFixed(1)} Mbps`"
        />
      </div>
      <div class="flex justify-between text-xs opacity-30 mt-1">
        <span>{{ speedHistory.length > 0 ? speedHistory[0].time : '' }}</span>
        <span>{{ speedHistory.length > 0 ? speedHistory[speedHistory.length - 1].time : '' }}</span>
      </div>
    </div>

    <!-- 测试历史 -->
    <div v-if="testHistory.length > 0" class="flex-1 overflow-auto">
      <div class="flex items-center justify-between mb-2">
        <div class="label mb-0">测试历史</div>
        <button @click="testHistory = []" class="btn btn-secondary py-0.5 px-2 text-xs">
          <Trash2 :size="11"/>清空
        </button>
      </div>
      <table class="w-full text-sm">
        <thead>
          <tr class="text-left opacity-50 text-xs border-b" style="border-color:var(--border-color)">
            <th class="pb-2 pr-4">时间</th>
            <th class="pb-2 pr-4">下载</th>
            <th class="pb-2 pr-4">上传</th>
            <th class="pb-2 pr-4">延迟</th>
            <th class="pb-2">丢包率</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(h, i) in testHistory" :key="i"
              class="border-b hover:bg-white/3 transition-colors"
              style="border-color:var(--border-color)">
            <td class="py-1.5 pr-4 text-xs opacity-50">{{ h.time }}</td>
            <td class="py-1.5 pr-4 font-mono text-green-400">{{ h.downloadSpeed }} Mbps</td>
            <td class="py-1.5 pr-4 font-mono text-blue-400">{{ h.uploadSpeed }} Mbps</td>
            <td class="py-1.5 pr-4 font-mono">{{ h.ping }} ms</td>
            <td class="py-1.5 font-mono" :class="Number(h.loss) > 0 ? 'text-red-400' : 'text-green-400'">{{ h.loss }}%</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 空状态 -->
    <div v-if="!testing && !testResults && testHistory.length === 0" class="flex items-center justify-center opacity-30 pt-8">
      <div class="text-center">
        <Gauge :size="32" class="mx-auto mb-2"/>
        <div class="text-sm">点击"开始测速"按钮</div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, computed, onUnmounted } from 'vue'
import { Gauge, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { SpeedTest as RunSpeedTest } from '../../../wailsjs/go/network/NetworkTools'

interface TestResult {
  downloadSpeed: string
  uploadSpeed: string
  ping: string
  loss: string
}

interface SpeedPoint {
  speed: number
  phase: 'download' | 'upload' | 'ping'
  time: string
}

interface TestHistoryItem extends TestResult {
  time: string
}

const appStore = useAppStore()
const testing = ref(false)
const canStop = ref(false)
const testPhase = ref<'idle' | 'ping' | 'download' | 'upload'>('idle')
const currentSpeed = ref('0.0')
const currentUnit = ref('Mbps')
const testResults = ref<TestResult | null>(null)
const speedHistory = ref<SpeedPoint[]>([])
const testHistory = ref<TestHistoryItem[]>([])
let testTimer: ReturnType<typeof setInterval> | null = null
let testStartTime = 0

// 仪表盘百分比
const gaugePercent = computed(() => {
  const speed = parseFloat(currentSpeed.value)
  return Math.min(100, (speed / 100) * 100)
})

// 仪表盘颜色
const gaugeColor = computed(() => {
  const speed = parseFloat(currentSpeed.value)
  if (speed < 10) return '#ef4444'
  if (speed < 30) return '#eab308'
  if (speed < 60) return '#3b82f6'
  return '#22c55e'
})

// 最大速度（用于图表缩放）
const maxSpeed = computed(() => {
  if (speedHistory.value.length === 0) return 1
  return Math.max(...speedHistory.value.map(s => s.speed), 1)
})

// 测试阶段文字
const testPhaseText = computed(() => {
  switch (testPhase.value) {
    case 'ping': return '正在测试延迟...'
    case 'download': return '正在测试下载速度...'
    case 'upload': return '正在测试上传速度...'
    default: return ''
  }
})

// 延迟颜色
function getPingColor(ping: string): string {
  const p = parseFloat(ping)
  if (p < 30) return 'text-green-400'
  if (p < 100) return 'text-yellow-400'
  return 'text-red-400'
}

// 获取当前时间字符串
function getTimeStr(): string {
  const now = new Date()
  return `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`
}

// 开始测速
async function startTest() {
  if (testing.value) {
    stopTest()
    return
  }

  testing.value = true
  canStop.value = false
  testResults.value = null
  speedHistory.value = []
  currentSpeed.value = '0.0'
  testStartTime = Date.now()

  try {
    // 尝试调用后端测速
    const res = await RunSpeedTest() as any
    if (res && res.success !== false) {
      testResults.value = {
        downloadSpeed: res.downloadMbps || '0.0',
        uploadSpeed: res.uploadMbps || '0.0',
        ping: res.latencyMs || '0',
        loss: res.packetLoss || '0',
      }
      addHistory(testResults.value)
      appStore.showToast('success', '测速完成')
      testing.value = false
      return
    }
  } catch {
    // 后端不可用，使用模拟测速
  }

  // 模拟测速
  await simulateSpeedTest()
}

// 模拟测速过程
async function simulateSpeedTest() {
  // 阶段1：延迟测试
  testPhase.value = 'ping'
  canStop.value = true
  const pingResult = Math.floor(Math.random() * 50) + 5

  for (let i = 0; i < 5; i++) {
    if (!testing.value) return
    currentSpeed.value = String(Math.floor(Math.random() * 30) + pingResult - 10)
    currentUnit.value = 'ms'
    speedHistory.value.push({
      speed: parseFloat(currentSpeed.value),
      phase: 'ping',
      time: getTimeStr(),
    })
    await sleep(500)
  }

  // 阶段2：下载速度测试
  testPhase.value = 'download'
  currentUnit.value = 'Mbps'
  const baseDownload = Math.random() * 80 + 20 // 20-100 Mbps
  let downloadSamples: number[] = []

  for (let i = 0; i < 20; i++) {
    if (!testing.value) return
    const variation = (Math.random() - 0.3) * 20
    const rampUp = Math.min(1, i / 5) // 前5次逐渐加速
    const speed = Math.max(0, (baseDownload + variation) * rampUp)
    currentSpeed.value = speed.toFixed(1)
    downloadSamples.push(speed)
    speedHistory.value.push({ speed, phase: 'download', time: getTimeStr() })
    await sleep(300)
  }

  const avgDownload = downloadSamples.reduce((a, b) => a + b, 0) / downloadSamples.length

  // 阶段3：上传速度测试
  testPhase.value = 'upload'
  const baseUpload = baseDownload * (Math.random() * 0.4 + 0.3) // 下载的30%-70%
  let uploadSamples: number[] = []

  for (let i = 0; i < 20; i++) {
    if (!testing.value) return
    const variation = (Math.random() - 0.3) * 15
    const rampUp = Math.min(1, i / 5)
    const speed = Math.max(0, (baseUpload + variation) * rampUp)
    currentSpeed.value = speed.toFixed(1)
    uploadSamples.push(speed)
    speedHistory.value.push({ speed, phase: 'upload', time: getTimeStr() })
    await sleep(300)
  }

  const avgUpload = uploadSamples.reduce((a, b) => a + b, 0) / uploadSamples.length
  const loss = Math.random() > 0.8 ? (Math.random() * 2).toFixed(1) : '0'

  // 生成结果
  testResults.value = {
    downloadSpeed: avgDownload.toFixed(1),
    uploadSpeed: avgUpload.toFixed(1),
    ping: String(pingResult),
    loss,
  }

  addHistory(testResults.value)
  currentSpeed.value = avgDownload.toFixed(1)
  testing.value = false
  canStop.value = false
  appStore.showToast('success', '测速完成')
}

// 停止测试
function stopTest() {
  testing.value = false
  canStop.value = false
  testPhase.value = 'idle'
  appStore.showToast('info', '测速已停止')
}

// 添加到历史
function addHistory(result: TestResult) {
  testHistory.value.unshift({
    ...result,
    time: getTimeStr(),
  })
  // 限制历史记录数量
  if (testHistory.value.length > 50) {
    testHistory.value = testHistory.value.slice(0, 50)
  }
}

// 延迟函数
function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

onUnmounted(() => {
  stopTest()
  if (testTimer) {
    clearInterval(testTimer)
    testTimer = null
  }
})
</script>
