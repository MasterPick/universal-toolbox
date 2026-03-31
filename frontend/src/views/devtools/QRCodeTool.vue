<template>
  <div class="page-container">
    <div class="page-title"><QrCode :size="20" class="text-primary-400"/>二维码工具</div>
    <div class="page-desc">生成二维码（需安装 qrcode.js）</div>
    <div class="card mb-4">
      <div class="label mb-2">输入内容</div>
      <input v-model="text" class="input-field mb-3" placeholder="输入文字、URL 等内容..."/>
      <button @click="generate" class="btn btn-primary"><Wand2 :size="14"/>生成二维码</button>
    </div>
    <div class="flex justify-center">
      <div id="qrcode-container" class="p-4 bg-white rounded-xl inline-block"></div>
    </div>
    <div v-if="!generated" class="text-center opacity-30 pt-8">
      <QrCode :size="48" class="mx-auto mb-2"/>
      <div class="text-sm">输入内容后点击生成</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { QrCode, Wand2 } from 'lucide-vue-next'
const text = ref(''), generated = ref(false)
let QRCode: any = null

onMounted(async () => {
  // 动态加载 qrcode 库
  try {
    const mod = await import('qrcode')
    QRCode = mod.default || mod
  } catch {}
})

async function generate() {
  if (!text.value || !QRCode) return
  const container = document.getElementById('qrcode-container')
  if (!container) return
  container.innerHTML = '<canvas></canvas>'
  const canvas = container.querySelector('canvas')
  await QRCode.toCanvas(canvas, text.value, { width: 256 })
  generated.value = true
}
</script>
