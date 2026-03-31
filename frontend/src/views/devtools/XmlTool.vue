<template>
  <div class="page-container">
    <div class="page-title"><Code2 :size="20" class="text-primary-400"/>XML 工具</div>
    <div class="page-desc">XML 格式化与美化</div>
    <div class="two-col flex-1 min-h-0">
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">输入 XML</div>
        <textarea v-model="input" class="textarea-field flex-1 min-h-0" placeholder="粘贴 XML 内容..." spellcheck="false"/>
        <button @click="format" class="btn btn-primary"><Wand2 :size="14"/>格式化</button>
      </div>
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label"><span>格式化结果</span>
          <button @click="copy" class="btn btn-secondary py-0.5 px-2 text-xs"><Copy :size="11"/>复制</button>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto">
          <span v-if="!output" class="opacity-30">格式化结果显示在这里...</span>
          <span v-else :class="isError ? 'text-red-400' : ''">{{ output }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Code2, Wand2, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { FormatXML } from '../../../wailsjs/go/devtools/DevTools'
const appStore = useAppStore()
const input = ref(''), output = ref(''), isError = ref(false)
async function format() {
  if (!input.value.trim()) return
  const res = await FormatXML(input.value)
  isError.value = !res.success
  output.value = res.success ? res.data : res.error
}
async function copy() {
  if (!output.value) return
  await navigator.clipboard.writeText(output.value)
  appStore.showToast('success', '已复制')
}
</script>
