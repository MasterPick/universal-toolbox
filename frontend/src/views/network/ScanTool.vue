<template>
  <div class="page-container">
    <div class="page-title"><Radar :size="20" class="text-primary-400"/>内网 IP 扫描</div>
    <div class="page-desc">扫描局域网内在线设备</div>
    <div class="card mb-4 flex gap-2 items-center">
      <span class="text-sm opacity-60">子网段：</span>
      <input v-model="subnet" class="input-field w-40" placeholder="如 192.168.1"/>
      <button @click="scan" class="btn btn-primary" :disabled="scanning">
        <Radar :size="14" :class="scanning?'loading-spin':''"/>
        {{ scanning ? `扫描中 ${progress}%...` : '开始扫描' }}
      </button>
    </div>
    <div v-if="results.length" class="flex-1 overflow-auto">
      <div class="text-xs opacity-50 mb-2">发现 {{ results.length }} 台在线设备</div>
      <div class="grid grid-cols-4 gap-2">
        <div v-for="r in results" :key="r.ip" class="card text-center">
          <div class="text-green-400 font-mono text-sm">{{ r.ip }}</div>
          <div class="text-xs opacity-50 mt-0.5">{{ r.latency }} ms</div>
        </div>
      </div>
    </div>
    <div v-else-if="!scanning" class="flex-1 flex items-center justify-center opacity-30">
      <div class="text-center"><Radar :size="48" class="mx-auto mb-2"/><div>点击开始扫描局域网</div></div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Radar } from 'lucide-vue-next'
import { ScanLAN, GetLocalSubnet } from '../../../wailsjs/go/network/NetworkTools'
const subnet = ref('192.168.1'), scanning = ref(false), progress = ref(0)
const results = ref<any[]>([])
async function scan() {
  scanning.value = true; results.value = []; progress.value = 0
  const timer = setInterval(() => { if (progress.value < 90) progress.value += 5 }, 500)
  try {
    results.value = (await ScanLAN(subnet.value) as any[]).sort((a,b)=>a.latency-b.latency)
  } finally {
    clearInterval(timer); progress.value = 100; scanning.value = false
  }
}
onMounted(async () => { subnet.value = await GetLocalSubnet() as string })
</script>
