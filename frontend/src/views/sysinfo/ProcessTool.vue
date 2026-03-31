<template>
  <div class="page-container">
    <div class="page-title"><Cpu :size="20" class="text-primary-400"/>进程管理</div>
    <div class="page-desc">查看并管理系统运行中的进程</div>
    <div class="flex gap-2 mb-4">
      <input v-model="filter" class="input-field flex-1" placeholder="按进程名过滤..."/>
      <button @click="loadProcs" class="btn btn-primary" :disabled="loading">
        <RefreshCw :size="14" :class="loading?'loading-spin':''"/>刷新
      </button>
    </div>
    <div class="flex-1 overflow-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-left opacity-50 text-xs border-b" style="border-color:var(--border-color)">
            <th class="pb-2 pr-4">进程名</th><th class="pb-2 pr-4">PID</th>
            <th class="pb-2 pr-4">CPU%</th><th class="pb-2 pr-4">内存%</th>
            <th class="pb-2">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in filteredProcs" :key="p.pid"
              class="border-b hover:bg-white/3 transition-colors"
              style="border-color:var(--border-color)">
            <td class="py-1.5 pr-4">{{ p.name }}</td>
            <td class="py-1.5 pr-4 font-mono text-xs opacity-60">{{ p.pid }}</td>
            <td class="py-1.5 pr-4">
              <span :class="p.cpu>50?'text-red-400':p.cpu>20?'text-yellow-400':'text-green-400'">
                {{ p.cpu.toFixed(1) }}%
              </span>
            </td>
            <td class="py-1.5 pr-4 opacity-70">{{ p.memory.toFixed(1) }}%</td>
            <td class="py-1.5">
              <button @click="killProc(p.pid)" class="btn btn-danger py-0.5 px-2 text-xs">
                <X :size="10"/>结束
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Cpu, RefreshCw, X } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetProcessList, KillProcess } from '../../../wailsjs/go/sysinfo/SysInfo'
const appStore = useAppStore()
const procs = ref<any[]>([]), filter = ref(''), loading = ref(false)
const filteredProcs = computed(() => filter.value
  ? procs.value.filter(p => p.name.toLowerCase().includes(filter.value.toLowerCase()))
  : procs.value.slice(0, 200)
)
async function loadProcs() {
  loading.value = true
  try { procs.value = (await GetProcessList() as any[] || []).sort((a,b)=>b.cpu-a.cpu) }
  finally { loading.value = false }
}
async function killProc(pid: number) {
  try { await KillProcess(pid); appStore.showToast('success', `进程 ${pid} 已终止`); loadProcs() }
  catch(e) { appStore.showToast('error', String(e)) }
}
onMounted(loadProcs)
</script>
