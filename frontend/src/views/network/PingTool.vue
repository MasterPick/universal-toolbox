<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Signal :size="20" class="text-primary-400"/>Ping 测试</div>
      <div class="page-desc">测试主机连通性与网络延迟</div>
    </div>

    <div class="card mb-4">
      <div class="flex gap-2">
        <input v-model="host" class="input-field flex-1"
          placeholder="输入主机名或 IP，如 google.com / 8.8.8.8"
          @keyup.enter="ping" />
        <button @click="ping" class="btn btn-primary" :disabled="pinging">
          <Signal :size="14" :class="pinging ? 'loading-spin' : ''"/>
          {{ pinging ? '测试中...' : '开始测试' }}
        </button>
        <button @click="clear" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
      </div>
    </div>

    <!-- 结果列表 -->
    <div class="flex-1 overflow-auto space-y-2">
      <div v-for="(r, i) in results" :key="i"
           :class="['card flex items-center gap-3', r.alive ? 'border-green-500/20' : 'border-red-500/20']">
        <div :class="['w-2 h-2 rounded-full shrink-0', r.alive ? 'bg-green-400' : 'bg-red-400']" />
        <div class="flex-1 font-mono text-sm">{{ r.host }}</div>
        <div v-if="r.alive" class="text-green-400 text-sm font-mono">{{ r.latencyMs.toFixed(1) }} ms</div>
        <div v-else class="text-red-400 text-sm">不可达</div>
        <div :class="['badge', r.alive ? 'badge-success' : 'badge-error']">
          {{ r.alive ? '在线' : '离线' }}
        </div>
      </div>
      <div v-if="results.length === 0" class="flex-1 flex items-center justify-center opacity-30 pt-16">
        <div class="text-center">
          <Signal :size="32" class="mx-auto mb-2"/>
          <div class="text-sm">输入主机名开始测试</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Signal, Trash2 } from 'lucide-vue-next'
import { PingHost } from '../../../wailsjs/go/network/NetworkTools'

interface PingResult { host: string; alive: boolean; latencyMs: number; error: string }

const host    = ref('')
const pinging = ref(false)
const results = ref<PingResult[]>([])

async function ping() {
  if (!host.value.trim()) return
  pinging.value = true
  const r = await PingHost(host.value.trim()) as PingResult
  results.value.unshift(r)
  pinging.value = false
}

function clear() { results.value = [] }
</script>
