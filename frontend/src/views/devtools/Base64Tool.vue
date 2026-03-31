<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Binary :size="20" class="text-primary-400"/>Base64 编解码</div>
      <div class="page-desc">将文本或文件进行 Base64 编码与解码</div>
    </div>

    <!-- 模式切换 -->
    <div class="tab-bar mb-4">
      <button :class="['tab-item', mode==='encode' && 'active']" @click="mode='encode'">编码</button>
      <button :class="['tab-item', mode==='decode' && 'active']" @click="mode='decode'">解码</button>
    </div>

    <div class="two-col flex-1 min-h-0">
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">{{ mode === 'encode' ? '原始文本' : 'Base64 字符串' }}</div>
        <textarea v-model="inputText" class="textarea-field flex-1 min-h-0"
          :placeholder="mode==='encode' ? '输入要编码的文本...' : '输入 Base64 字符串...'"
          spellcheck="false" />
        <button @click="process" class="btn btn-primary w-full">
          <Wand2 :size="14"/>{{ mode === 'encode' ? '编码' : '解码' }}
        </button>
      </div>
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>{{ mode === 'encode' ? 'Base64 结果' : '解码结果' }}</span>
          <button @click="copyOutput" class="btn btn-secondary py-0.5 px-2 text-xs">
            <Copy :size="11"/>复制
          </button>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto">
          <span v-if="!output" class="opacity-30">结果将显示在这里...</span>
          <span v-else :class="isError ? 'text-red-400' : ''">{{ output }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Binary, Wand2, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { Base64Encode, Base64Decode } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const mode = ref<'encode' | 'decode'>('encode')
const inputText = ref('')
const output = ref('')
const isError = ref(false)

async function process() {
  if (!inputText.value.trim()) return
  isError.value = false
  if (mode.value === 'encode') {
    const res = await Base64Encode(inputText.value)
    output.value = res.data
  } else {
    const res = await Base64Decode(inputText.value)
    if (res.success) {
      output.value = res.data
    } else {
      output.value = res.error
      isError.value = true
      appStore.showToast('error', res.error)
    }
  }
}

async function copyOutput() {
  if (!output.value) return
  await navigator.clipboard.writeText(output.value)
  appStore.showToast('success', '已复制到剪贴板')
}
</script>
