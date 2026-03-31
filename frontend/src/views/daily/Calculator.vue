<template>
  <div class="page-container">
    <div class="page-title"><Calculator :size="20" class="text-primary-400"/>计算器</div>
    <div class="tab-bar mb-4">
      <button :class="['tab-item', mode==='standard'&&'active']" @click="mode='standard'">标准</button>
      <button :class="['tab-item', mode==='scientific'&&'active']" @click="mode='scientific'">科学</button>
    </div>
    <div class="calc-screen card mb-4">
      <div class="text-xs opacity-40 h-4">{{ expression }}</div>
      <div class="text-3xl font-mono font-light mt-1 truncate">{{ display }}</div>
    </div>
    <div class="calc-grid">
      <template v-if="mode==='scientific'">
        <button v-for="b in sciButtons" :key="b" @click="pressSci(b)" class="calc-btn calc-btn-fn">{{ b }}</button>
      </template>
      <button @click="clear"           class="calc-btn calc-btn-clear col-span-2">AC</button>
      <button @click="press('%')"      class="calc-btn calc-btn-fn">%</button>
      <button @click="press('/')"      class="calc-btn calc-btn-op">÷</button>
      <button v-for="n in [7,8,9]"  :key="n" @click="press(String(n))" class="calc-btn">{{ n }}</button>
      <button @click="press('*')"      class="calc-btn calc-btn-op">×</button>
      <button v-for="n in [4,5,6]"  :key="n" @click="press(String(n))" class="calc-btn">{{ n }}</button>
      <button @click="press('-')"      class="calc-btn calc-btn-op">−</button>
      <button v-for="n in [1,2,3]"  :key="n" @click="press(String(n))" class="calc-btn">{{ n }}</button>
      <button @click="press('+')"      class="calc-btn calc-btn-op">+</button>
      <button @click="press('0')"      class="calc-btn col-span-2">0</button>
      <button @click="press('.')"      class="calc-btn">.</button>
      <button @click="calculate"       class="calc-btn calc-btn-eq">=</button>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Calculator } from 'lucide-vue-next'
const mode = ref<'standard'|'scientific'>('standard')
const display = ref('0'), expression = ref(''), current = ref('')
const sciButtons = ['sin','cos','tan','√','x²','ln','log','π']

function press(val: string) {
  if (['+','-','*','/','%'].includes(val)) {
    expression.value = display.value + ' ' + val
    current.value = ''
  } else {
    current.value += val
    display.value = current.value
  }
}

function calculate() {
  try {
    const expr = expression.value + ' ' + display.value
    // 简单四则运算（前端直接计算）
    const result = Function('"use strict"; return (' + expr.replace('÷','/').replace('×','*') + ')')()
    display.value = String(parseFloat(result.toFixed(10)))
    expression.value = expr + ' ='
    current.value = display.value
  } catch { display.value = '错误' }
}

function pressSci(fn: string) {
  const v = parseFloat(display.value)
  const map: Record<string,number> = {
    'sin': Math.sin(v*Math.PI/180), 'cos': Math.cos(v*Math.PI/180),
    'tan': Math.tan(v*Math.PI/180), '√': Math.sqrt(v),
    'x²': v*v, 'ln': Math.log(v), 'log': Math.log10(v), 'π': Math.PI
  }
  display.value = String(parseFloat((map[fn]??0).toFixed(10)))
}

function clear() { display.value='0'; expression.value=''; current.value='' }
</script>
<style scoped>
.calc-screen { padding: 12px 16px; }
.calc-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.calc-btn {
  background: var(--bg-hover); border: 1px solid var(--border-color);
  border-radius: 10px; padding: 14px 8px; font-size: 16px; cursor: pointer;
  color: var(--text-primary); transition: all 0.1s; text-align: center;
}
.calc-btn:hover  { background: rgba(255,255,255,0.1); }
.calc-btn:active { transform: scale(0.95); }
.calc-btn-op    { color: #818cf8; font-size: 18px; }
.calc-btn-eq    { background: #6366f1; color: #fff; }
.calc-btn-eq:hover { background: #818cf8; }
.calc-btn-clear { background: rgba(239,68,68,0.15); color: #f87171; }
.calc-btn-fn    { color: #94a3b8; font-size: 12px; }
</style>
