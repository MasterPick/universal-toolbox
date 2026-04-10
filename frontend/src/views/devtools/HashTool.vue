<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Hash :size="20" class="text-primary-400"/>哈希计算</div>
      <div class="page-desc">计算文本的 MD5、SHA1、SHA256、SHA512 哈希值</div>
    </div>

    <div class="card mb-4">
      <div class="label">输入文本</div>
      <textarea v-model="inputText" class="textarea-field" rows="4"
        placeholder="输入要计算哈希的文本..." spellcheck="false" />
      <div class="flex gap-2 mt-3">
        <button @click="calcAll" class="btn btn-primary" :disabled="calculating">
          <Wand2 :size="14"/>{{ calculating ? '计算中...' : '全部计算' }}
        </button>
        <button @click="clearAll" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
      </div>
    </div>

    <div class="space-y-3 flex-1 overflow-auto">
      <div v-for="item in results" :key="item.algo" class="card">
        <div class="flex items-center justify-between mb-2">
          <span class="font-mono font-bold text-sm text-primary-400">{{ item.algo }}</span>
          <div class="flex items-center gap-2">
            <span v-if="item.error" class="text-xs text-red-400">{{ item.error }}</span>
            <button v-if="item.value" @click="copyHash(item.value)" class="btn btn-secondary py-0.5 px-2 text-xs">
              <Copy :size="11"/>复制
            </button>
          </div>
        </div>
        <div class="code-output text-sm break-all">
          <span v-if="!item.value && !item.error" class="opacity-30">点击"全部计算"获取结果</span>
          <span v-else-if="item.error" class="text-red-400">{{ item.error }}</span>
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
import { CalcMD5, CalcSHA1, CalcSHA256, CalcSHA512 } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const inputText = ref('')
const calculating = ref(false)
const results = ref([
  { algo: 'MD5',    value: '', error: '' },
  { algo: 'SHA1',   value: '', error: '' },
  { algo: 'SHA256', value: '', error: '' },
  { algo: 'SHA512', value: '', error: '' },
])

async function calcAll() {
  if (!inputText.value) return
  calculating.value = true
  // 清空之前的错误
  results.value.forEach(r => { r.value = ''; r.error = '' })
  try {
    const [md5, sha1, sha256, sha512] = await Promise.all([
      CalcMD5(inputText.value).catch((e: unknown) => ({ success: false, error: String(e) })),
      CalcSHA1(inputText.value).catch((e: unknown) => ({ success: false, error: String(e) })),
      CalcSHA256(inputText.value).catch((e: unknown) => ({ success: false, error: String(e) })),
      CalcSHA512(inputText.value).catch((e: unknown) => ({ success: false, error: String(e) })),
    ])
    if (md5.success !== false) { results.value[0].value = md5.data } else { results.value[0].error = md5.error }
    if (sha1.success !== false) { results.value[1].value = sha1.data } else { results.value[1].error = sha1.error }
    if (sha256.success !== false) { results.value[2].value = sha256.data } else { results.value[2].error = sha256.error }
    if (sha512.success !== false) { results.value[3].value = sha512.data } else { results.value[3].error = sha512.error }

    appStore.showToast('success', '哈希计算完成')
  } catch (e) {
    appStore.showToast('error', '哈希计算失败: ' + String(e))
  } finally {
    calculating.value = false
  }
}

async function copyHash(val: string) {
  if (!val) return
  try {
    await navigator.clipboard.writeText(val)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

function clearAll() {
  inputText.value = ''
  results.value.forEach(r => { r.value = ''; r.error = '' })
}
</script>
