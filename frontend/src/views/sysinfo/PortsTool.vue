<template>
  <div class="page-container">
    <div class="page-title"><Network :size="20" class="text-primary-400"/>端口查看</div>
    <div class="page-desc">查看端口占用，一键释放</div>
    <div class="flex gap-2 mb-4">
      <input v-model="filterPort" class="input-field w-40" placeholder="过滤端口号..." type="number"/>
      <button @click="loadPorts" class="btn btn-primary" :disabled="loading">
        <RefreshCw :size="14" :class="loading?'loading-spin':''"/>刷新
      </button>
    </div>
    <div class="flex-1 overflow-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-left opacity-50 text-xs border-b" style="border-color:var(--border-color)">
            <th class="pb-2 pr-4">端口</th><th class="pb-2 pr-4">进程</th>
            <th class="pb-2 pr-4">PID</th><th class="pb-2 pr-4">状态</th>
            <th class="pb-2">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in filteredPorts" :key="p.port+'-'+p.pid"
              class="border-b hover:bg-white/3 transition-colors"
              style="border-color:var(--border-color)">
            <td class="py-2 pr-4 font-mono text-primary-400">{{ p.port }}</td>
            <td class="py-2 pr-4">{{ p.process || '-' }}</td>
            <td class="py-2 pr-4 font-mono text-xs opacity-60">{{ p.pid }}</td>
            <td class="py-2 pr-4">
              <span :class="['badge', p.status==='LISTEN'?'badge-success':'badge-info']">{{ p.status }}</span>
            </td>
            <td class="py-2">
              <button @click="releasePort(p.port)" class="btn btn-danger py-0.5 px-2 text-xs">
                <Zap :size="10"/>释放
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="filteredPorts.length===0 && !loading" class="text-center opacity-30 pt-16">
        <Network :size="36" class="mx-auto mb-2"/><div>暂无数据，点击刷新</div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Network, RefreshCw, Zap } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetPortList, ReleasePort } from '../../../wailsjs/go/sysinfo/SysInfo'
const appStore = useAppStore()
const ports = ref<any[]>([])
const filterPort = ref<number|null>(null)
const loading = ref(false)
const filteredPorts = computed(() =>
  filterPort.value ? ports.value.filter(p => p.port === filterPort.value) : ports.value
)
async function loadPorts() {
  loading.value = true
  try { ports.value = await GetPortList() as any[] } finally { loading.value = false }
}
async function releasePort(port: number) {
  try {
    await ReleasePort(port)
    appStore.showToast('success', `端口 ${port} 已释放`)
    loadPorts()
  } catch(e) { appStore.showToast('error', String(e)) }
}
onMounted(loadPorts)
</script>
