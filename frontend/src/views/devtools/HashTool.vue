<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Hash :size="20" class="text-primary-400"/>哈希计算</div>
      <div class="page-desc">计算文本的 MD5、SHA1、SHA256 哈希值</div>
    </div>

    <div class="card mb-4">
      <div class="label">输入文本</div>
      <textarea v-model="inputText" class="textarea-field" rows="4"
        placeholder="输入要计算哈希的文本..." spellcheck="false" />
      <div class="flex gap-2 mt-3">
        <button @click="calcAll" class="btn btn-primary"><Wand2 :size="14"/>全部计算</button>
        <button @click="clearAll" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
      </div>
    </div>

    <div class="space-y-3 flex-1 overflow-auto">
      <div v-for="item in results" :key="item.algo" class="card">
        <div class="flex items-center justify-between mb-2">
          <span class="font-mono font-bold text-sm text-primary-400">{{ item.algo }}</span>
          <button @click="copyHash(item.value)" class="btn btn-secondary py-0.5 px-2 text-xs">
            <Copy :size="11"/>复制
          </button>
        </div>
        <div class="code-output text-sm break-all">
          <span v-if="!item.value" class="opacity-30">点击"全部计算"获取结果</span>
          <span v-else>{{ item.value }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Hash, Wand2, Trash2, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { CalcMD5, CalcSHA1, CalcSHA256 } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const inputText = ref('')
const results = ref([
  { algo: 'MD5',    value: '' },
  { algo: 'SHA1',   value: '' },
  { algo: 'SHA256', value: '' },
])

async function calcAll() {
  if (!inputText.value) return
  const [md5, sha1, sha256] = await Promise.all([
    CalcMD5(inputText.value),
    CalcSHA1(inputText.value),
    CalcSHA256(inputText.value),
  ])
  results.value[0].value = md5.data
  results.value[1].value = sha1.data
  results.value[2].value = sha256.data
  appStore.showToast('success', '哈希计算完成')
}

async function copyHash(val: string) {
  if (!val) return
  await navigator.clipboard.writeText(val)
  appStore.showToast('success', '已复制到剪贴板')
}

function clearAll() {
  inputText.value = ''
  results.value.forEach(r => r.value = '')
}
</script>
