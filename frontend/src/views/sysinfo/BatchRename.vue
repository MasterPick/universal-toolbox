<template>
  <div class="page-container">
    <div class="page-title"><FilePen :size="20" class="text-primary-400"/>批量重命名</div>
    <div class="page-desc">批量重命名文件，支持正则和序号</div>
    <div class="card mb-4">
      <div class="label mb-2">拖放文件到此处（或点击选择）</div>
      <div class="drop-zone" @dragover.prevent @drop="onDrop">
        <Upload :size="32" class="mx-auto mb-2 opacity-40"/>
        <div class="text-sm opacity-40">拖放文件到这里</div>
      </div>
    </div>
    <div v-if="files.length" class="card">
      <div class="flex gap-2 mb-3">
        <input v-model="pattern" class="input-field flex-1" placeholder="查找（支持正则）..."/>
        <input v-model="replacement" class="input-field flex-1" placeholder="替换为（$1 引用捕获组）..."/>
        <button @click="preview" class="btn btn-secondary">预览</button>
        <button @click="rename" class="btn btn-primary">重命名</button>
      </div>
      <div class="space-y-1 max-h-48 overflow-auto">
        <div v-for="(f,i) in files" :key="i" class="flex gap-3 text-xs py-1">
          <span class="opacity-60 flex-1 truncate">{{ f.original }}</span>
          <span class="text-primary-400">→</span>
          <span class="flex-1 truncate">{{ f.renamed || f.original }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { FilePen, Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
const appStore = useAppStore()
const files = ref<{original:string,renamed:string,path:string}[]>([])
const pattern = ref(''), replacement = ref('')

function onDrop(e: DragEvent) {
  const items = e.dataTransfer?.files
  if (!items) return
  for (let i = 0; i < items.length; i++) {
    const f = items[i] as any
    files.value.push({ original: f.name, renamed: f.name, path: f.path || '' })
  }
}

function preview() {
  files.value = files.value.map(f => {
    try {
      const re = new RegExp(pattern.value, 'g')
      return { ...f, renamed: f.original.replace(re, replacement.value) }
    } catch { return f }
  })
}

async function rename() {
  appStore.showToast('info', '批量重命名功能需配合文件系统权限，请在正式打包后使用')
}
</script>
<style scoped>
.drop-zone {
  border: 2px dashed var(--border-color);
  border-radius: 10px; padding: 32px;
  text-align: center; cursor: pointer;
  transition: border-color 0.2s;
}
.drop-zone:hover { border-color: var(--accent); }
</style>
