<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Fingerprint :size="20" class="text-primary-400"/>UUID 生成器</div>
      <div class="page-desc">批量生成标准 UUID v4</div>
    </div>
    <div class="card mb-4 flex items-center gap-4">
      <div class="label mb-0">生成数量</div>
      <input v-model.number="count" type="number" min="1" max="100"
        class="input-field w-24" />
      <button @click="generate" class="btn btn-primary"><RefreshCw :size="14"/>生成</button>
      <button @click="copyAll"  class="btn btn-secondary"><Copy :size="14"/>复制全部</button>
      <button @click="clear"    class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
    </div>
    <div class="code-output flex-1 min-h-0 overflow-auto">
      <span v-if="!uuids" class="opacity-30">生成的 UUID 将显示在这里...</span>
      <span v-else>{{ uuids }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Fingerprint, RefreshCw, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GenerateUUIDs } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const count = ref(10)
const uuids = ref('')

async function generate() {
  const res = await GenerateUUIDs(count.value)
  uuids.value = res.data
  appStore.showToast('success', `已生成 ${count.value} 个 UUID`)
}

async function copyAll() {
  if (!uuids.value) return
  await navigator.clipboard.writeText(uuids.value)
  appStore.showToast('success', '已复制到剪贴板')
}

function clear() { uuids.value = '' }
</script>
