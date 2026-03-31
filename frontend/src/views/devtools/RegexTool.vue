<template>
  <div class="page-container">
    <div class="page-title"><Regex :size="20" class="text-primary-400"/>正则测试</div>
    <div class="page-desc">测试正则表达式匹配结果</div>
    <div class="card mb-3 flex gap-2 items-center">
      <span class="opacity-50 font-mono">/</span>
      <input v-model="pattern" class="input-field flex-1" placeholder="输入正则表达式，如 \d{4}-\d{2}-\d{2}" />
      <span class="opacity-50 font-mono">/g</span>
      <button @click="test" class="btn btn-primary"><Play :size="14"/>测试</button>
    </div>
    <div class="label mb-1">测试文本</div>
    <textarea v-model="testText" class="textarea-field mb-3" rows="5" placeholder="输入要测试的文本..." />
    <div class="label mb-1">匹配结果</div>
    <div class="code-output flex-1 overflow-auto">{{ result || '匹配结果显示在这里...' }}</div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Regex, Play } from 'lucide-vue-next'
import { RegexTest } from '../../../wailsjs/go/devtools/DevTools'
const pattern = ref(''), testText = ref(''), result = ref('')
async function test() {
  if (!pattern.value || !testText.value) return
  const res = await RegexTest(pattern.value, testText.value)
  result.value = res.success ? res.data : res.error
}
</script>
