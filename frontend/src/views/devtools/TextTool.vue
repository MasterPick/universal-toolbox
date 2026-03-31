<template>
  <div class="page-container">
    <div class="page-title"><FileText :size="20" class="text-primary-400"/>文本工具</div>
    <div class="page-desc">文本对比 · 查找替换 · 字符统计</div>
    <div class="tab-bar mb-4">
      <button v-for="t in tabs" :key="t.id" :class="['tab-item', tab===t.id&&'active']" @click="tab=t.id">{{ t.label }}</button>
    </div>
    <!-- 字符统计 -->
    <template v-if="tab==='stats'">
      <textarea v-model="statsInput" class="textarea-field mb-3" rows="6" placeholder="输入文本进行统计..." />
      <button @click="calcStats" class="btn btn-primary mb-3"><BarChart2 :size="14"/>统计</button>
      <div class="code-output">{{ statsResult || '统计结果显示在这里...' }}</div>
    </template>
    <!-- 查找替换 -->
    <template v-if="tab==='replace'">
      <div class="grid grid-cols-2 gap-2 mb-2">
        <input v-model="search" class="input-field" placeholder="查找内容..." />
        <input v-model="replace" class="input-field" placeholder="替换为..." />
      </div>
      <div class="flex items-center gap-2 mb-2">
        <label class="flex items-center gap-1 text-sm cursor-pointer">
          <input type="checkbox" v-model="useRegex" class="mr-1"/>使用正则
        </label>
        <button @click="doReplace" class="btn btn-primary"><Replace :size="14"/>替换</button>
      </div>
      <textarea v-model="replaceInput" class="textarea-field flex-1 min-h-0" rows="8" placeholder="输入要处理的文本..." />
    </template>
    <!-- 文本对比 -->
    <template v-if="tab==='compare'">
      <div class="two-col mb-3" style="height: 200px">
        <textarea v-model="text1" class="textarea-field h-full" placeholder="原文..." />
        <textarea v-model="text2" class="textarea-field h-full" placeholder="新文..." />
      </div>
      <button @click="compare" class="btn btn-primary mb-3"><GitCompare :size="14"/>对比</button>
      <div class="code-output flex-1 overflow-auto">{{ compareResult || '对比结果显示在这里...' }}</div>
    </template>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { FileText, BarChart2, Replace, GitCompare } from 'lucide-vue-next'
import { TextStats, TextReplace, TextCompare } from '../../../wailsjs/go/devtools/DevTools'
const tab = ref('stats')
const tabs = [{ id:'stats', label:'字符统计' }, { id:'replace', label:'查找替换' }, { id:'compare', label:'文本对比' }]
const statsInput = ref(''), statsResult = ref('')
const search = ref(''), replace = ref(''), replaceInput = ref(''), useRegex = ref(false)
const text1 = ref(''), text2 = ref(''), compareResult = ref('')
async function calcStats() {
  const res = await TextStats(statsInput.value)
  statsResult.value = res.data
}
async function doReplace() {
  const res = await TextReplace(replaceInput.value, search.value, replace.value, useRegex.value)
  replaceInput.value = res.success ? res.data : res.error
}
async function compare() {
  const res = await TextCompare(text1.value, text2.value)
  compareResult.value = res.data
}
</script>
