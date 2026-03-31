<template>
  <div class="page-container">
    <div class="page-title"><ArrowLeftRight :size="20" class="text-primary-400"/>单位换算</div>
    <div class="page-desc">长度 · 重量 · 温度 · 速度</div>
    <div class="tab-bar mb-4">
      <button v-for="t in types" :key="t.id" :class="['tab-item', type===t.id&&'active']" @click="type=t.id">{{ t.label }}</button>
    </div>
    <div class="card">
      <div class="flex gap-3 items-end mb-4">
        <div class="flex-1">
          <div class="label mb-1">数值</div>
          <input v-model.number="value" type="number" class="input-field" placeholder="输入数值..." />
        </div>
        <div class="flex-1">
          <div class="label mb-1">从</div>
          <select v-model="fromUnit" class="input-field">
            <option v-for="u in currentUnits" :key="u.id" :value="u.id">{{ u.label }}</option>
          </select>
        </div>
        <div class="flex-1">
          <div class="label mb-1">转换到</div>
          <select v-model="toUnit" class="input-field">
            <option v-for="u in currentUnits" :key="u.id" :value="u.id">{{ u.label }}</option>
          </select>
        </div>
        <button @click="convert" class="btn btn-primary"><ArrowLeftRight :size="14"/>转换</button>
      </div>
      <div v-if="result" class="code-output text-center text-lg">{{ result }}</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { ArrowLeftRight } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { ConvertLength, ConvertWeight, ConvertTemperature, ConvertSpeed } from '../../../wailsjs/go/daily/DailyTools'
const appStore = useAppStore()
const type = ref('length'), value = ref(1), fromUnit = ref('m'), toUnit = ref('km'), result = ref('')
const types = [{id:'length',label:'长度'},{id:'weight',label:'重量'},{id:'temperature',label:'温度'},{id:'speed',label:'速度'}]
const unitMap: Record<string, {id:string,label:string}[]> = {
  length: [{id:'mm',label:'毫米'},{id:'cm',label:'厘米'},{id:'m',label:'米'},{id:'km',label:'千米'},{id:'inch',label:'英寸'},{id:'foot',label:'英尺'},{id:'mile',label:'英里'}],
  weight: [{id:'mg',label:'毫克'},{id:'g',label:'克'},{id:'kg',label:'千克'},{id:'t',label:'吨'},{id:'oz',label:'盎司'},{id:'lb',label:'磅'}],
  temperature: [{id:'C',label:'摄氏度 °C'},{id:'F',label:'华氏度 °F'},{id:'K',label:'开尔文 K'}],
  speed: [{id:'ms',label:'米/秒'},{id:'kmh',label:'千米/时'},{id:'mph',label:'英里/时'},{id:'knot',label:'节'}],
}
const currentUnits = computed(() => unitMap[type.value] || [])
async function convert() {
  if (!value.value) return
  try {
    let res: any
    if (type.value === 'length')      res = await ConvertLength(value.value, fromUnit.value, toUnit.value)
    else if (type.value === 'weight') res = await ConvertWeight(value.value, fromUnit.value, toUnit.value)
    else if (type.value === 'temperature') res = await ConvertTemperature(value.value, fromUnit.value, toUnit.value)
    else res = await ConvertSpeed(value.value, fromUnit.value, toUnit.value)
    result.value = res.formula
  } catch(e) { appStore.showToast('error', String(e)) }
}
</script>
