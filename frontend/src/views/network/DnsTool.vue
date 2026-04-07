<template>
  <!-- DNS 查询工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Globe2 :size="20" class="text-primary-400" />
        DNS 查询
      </div>
      <div class="page-desc">域名解析 · DNS 记录查询 · 支持国内外 DNS 服务器</div>
    </div>

    <!-- 查询输入 -->
    <div class="card mb-4">
      <div class="flex gap-3 flex-wrap">
        <input
          v-model="domain"
          class="input-field flex-1 min-w-[200px]"
          placeholder="输入域名，如 example.com"
          @keyup.enter="query"
        />
        <select v-model="recordType" class="input-field w-28">
          <option value="A">A 记录</option>
          <option value="AAAA">AAAA 记录</option>
          <option value="CNAME">CNAME</option>
          <option value="MX">MX 记录</option>
          <option value="TXT">TXT 记录</option>
          <option value="NS">NS 记录</option>
          <option value="SOA">SOA</option>
        </select>
        <select v-model="dnsServer" class="input-field w-36">
          <optgroup label="国内 DNS">
            <option value="alipay">阿里 DNS</option>
            <option value="tencent">腾讯 DNS</option>
            <option value="dnspod">DNSPod</option>
          </optgroup>
          <optgroup label="国际 DNS">
            <option value="google">Google DNS</option>
            <option value="cloudflare">Cloudflare</option>
          </optgroup>
        </select>
        <button @click="query" :disabled="isLoading" class="btn btn-primary">
          <Loader2 v-if="isLoading" :size="14" class="loading-spin" />
          <Search v-else :size="14" />
          查询
        </button>
      </div>
    </div>

    <!-- 快捷域名 -->
    <div class="flex gap-2 mb-4 flex-wrap">
      <span class="text-xs opacity-50">快捷：</span>
      <button
        v-for="d in quickDomains"
        :key="d"
        @click="domain = d; query()"
        class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10"
      >
        {{ d }}
      </button>
    </div>

    <!-- 查询结果 -->
    <div class="flex-1 flex flex-col min-h-0">
      <!-- 空状态 -->
      <div v-if="!result && !isLoading" class="flex-1 flex items-center justify-center opacity-30">
        <div class="text-center">
          <Globe2 :size="40" class="mx-auto mb-2 opacity-50" />
          <div class="text-sm">输入域名进行 DNS 查询</div>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-else-if="isLoading" class="flex-1 flex items-center justify-center">
        <div class="text-center">
          <Loader2 :size="32" class="mx-auto mb-2 loading-spin text-primary-400" />
          <div class="text-sm opacity-60">正在查询...</div>
        </div>
      </div>

      <!-- 结果展示 -->
      <div v-else-if="result" class="flex-1 overflow-auto">
        <!-- 错误 -->
        <div v-if="result.error" class="card bg-red-500/10 border-red-500/30">
          <div class="flex items-center gap-2 text-red-400">
            <AlertCircle :size="16" />
            <span>{{ result.error }}</span>
          </div>
        </div>

        <!-- 成功结果 -->
        <div v-else class="space-y-3">
          <!-- 基本信息 -->
          <div class="card">
            <div class="flex items-center gap-2 mb-3">
              <CheckCircle :size="16" class="text-green-400" />
              <span class="font-semibold">{{ result.domain }}</span>
              <span class="badge badge-info">{{ result.type }} 记录</span>
            </div>
            <div class="text-xs opacity-50">
              查询时间: {{ result.queryTime }}ms · DNS 服务器: {{ result.dnsServer || '系统默认' }}
            </div>
          </div>

          <!-- 记录列表 -->
          <div class="card">
            <div class="font-semibold mb-3">解析结果 ({{ result.records?.length || 0 }} 条)</div>
            <div class="space-y-2">
              <div
                v-for="(record, idx) in result.records"
                :key="idx"
                class="flex items-center gap-3 p-3 rounded-lg bg-white/5"
              >
                <div class="w-8 h-8 rounded-full flex items-center justify-center" :class="getRecordColor(result.type)">
                  {{ idx + 1 }}
                </div>
                <div class="flex-1 font-mono text-sm">
                  {{ record.value || record }}
                </div>
                <div v-if="record.ttl" class="text-xs opacity-50">
                  TTL: {{ record.ttl }}
                </div>
                <button @click="copyValue(record.value || record)" class="btn btn-secondary text-xs py-1">
                  <Copy :size="12" />
                </button>
              </div>
            </div>
          </div>

          <!-- 原始响应 -->
          <div class="card">
            <div class="flex items-center justify-between mb-2">
              <div class="font-semibold text-sm">原始响应</div>
              <button @click="showRaw = !showRaw" class="text-xs opacity-50">
                {{ showRaw ? '收起' : '展开' }}
              </button>
            </div>
            <div v-if="showRaw" class="code-output text-xs">
              <pre>{{ JSON.stringify(result.raw, null, 2) }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 查询历史 -->
    <div v-if="history.length" class="mt-4">
      <div class="flex items-center justify-between mb-2">
        <div class="text-sm opacity-60">查询历史</div>
        <button @click="history = []" class="text-xs opacity-50">清空</button>
      </div>
      <div class="flex gap-2 flex-wrap">
        <button
          v-for="h in history.slice(0, 10)"
          :key="h"
          @click="domain = h; query()"
          class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10 font-mono"
        >
          {{ h }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Globe2, Search, Loader2, AlertCircle, CheckCircle, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const domain = ref('')
const recordType = ref('A')
const dnsServer = ref('alipay') // 默认使用阿里 DNS
const isLoading = ref(false)
const showRaw = ref(false)
const result = ref<any>(null)
const history = ref<string[]>([])

// DNS 服务器配置
const DNS_SERVERS: Record<string, { name: string; dohUrl: string }> = {
  alipay: {
    name: '阿里 DNS (223.5.5.5)',
    dohUrl: 'https://dns.alidns.com/dns-query'
  },
  tencent: {
    name: '腾讯 DNS (119.29.29.29)',
    dohUrl: 'https://doh.pub/dns-query'
  },
  dnspod: {
    name: 'DNSPod (119.29.29.29)',
    dohUrl: 'https://sm2.doh.pub/dns-query'
  },
  google: {
    name: 'Google DNS (8.8.8.8)',
    dohUrl: 'https://dns.google/resolve'
  },
  cloudflare: {
    name: 'Cloudflare (1.1.1.1)',
    dohUrl: 'https://cloudflare-dns.com/dns-query'
  }
}

// 快捷域名
const quickDomains = ['google.com', 'github.com', 'baidu.com', 'qq.com', 'openai.com']

// 查询 DNS
async function query() {
  if (!domain.value.trim()) return

  isLoading.value = true
  result.value = null

  const startTime = Date.now()

  try {
    const server = DNS_SERVERS[dnsServer.value]
    const data = await queryDoh(server.dohUrl, domain.value, recordType.value)

    if (data.Status !== 0) {
      result.value = {
        error: getDnsStatus(data.Status),
        domain: domain.value,
        type: recordType.value
      }
    } else {
      const records = data.Answer?.map((a: any) => ({
        value: a.data,
        ttl: a.TTL
      })) || []

      result.value = {
        domain: domain.value,
        type: recordType.value,
        records,
        queryTime: Date.now() - startTime,
        dnsServer: server.name,
        raw: data
      }

      // 添加到历史
      if (!history.value.includes(domain.value)) {
        history.value.unshift(domain.value)
        if (history.value.length > 20) history.value.pop()
        localStorage.setItem('dnsHistory', JSON.stringify(history.value))
      }
    }
  } catch (err) {
    result.value = {
      error: '查询失败: ' + String(err),
      domain: domain.value,
      type: recordType.value
    }
  } finally {
    isLoading.value = false
  }
}

// DoH 查询
async function queryDoh(dohUrl: string, domain: string, type: string): Promise<any> {
  // Google DoH 使用 GET 请求，参数在 URL 中
  if (dohUrl.includes('dns.google')) {
    const url = `${dohUrl}?name=${domain}&type=${type}`
    const response = await fetch(url)
    return await response.json()
  }

  // 阿里/腾讯/Cloudflare DoH 使用 POST 请求，DNS wire format
  const dnsPacket = createDnsQuery(domain, type)

  const response = await fetch(dohUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/dns-message',
      'Accept': 'application/dns-message'
    },
    body: dnsPacket
  })

  // 解析 DNS 响应
  const responseBuffer = await response.arrayBuffer()
  return parseDnsResponse(responseBuffer)
}

// 创建 DNS 查询包 (wire format)
function createDnsQuery(domain: string, type: string): ArrayBuffer {
  const typeMap: Record<string, number> = {
    A: 1, AAAA: 28, CNAME: 5, MX: 15, TXT: 16, NS: 2, SOA: 6
  }
  const qtype = typeMap[type] || 1

  // DNS header (12 bytes) + Question
  const labels = domain.split('.').map(label => {
    const buf = new Uint8Array(label.length + 1)
    buf[0] = label.length
    for (let i = 0; i < label.length; i++) buf[i + 1] = label.charCodeAt(i)
    return buf
  })

  const header = new Uint8Array([
    0x00, 0x01, // Transaction ID
    0x01, 0x00, // Flags: standard query
    0x00, 0x01, // Questions: 1
    0x00, 0x00, // Answers: 0
    0x00, 0x00, // Authority: 0
    0x00, 0x00  // Additional: 0
  ])

  const question = new Uint8Array([
    ...labels.flat(),
    0x00, // Root label
    (qtype >> 8) & 0xff, qtype & 0xff, // QTYPE
    0x00, 0x01 // QCLASS: IN
  ])

  const packet = new Uint8Array(header.length + question.length)
  packet.set(header, 0)
  packet.set(question, header.length)

  return packet.buffer
}

// 解析 DNS 响应包
function parseDnsResponse(buffer: ArrayBuffer): any {
  const view = new DataView(buffer)
  const uint8 = new Uint8Array(buffer)

  // 解析 header
  const flags = view.getUint16(2)
  const status = flags & 0x000F
  const qdcount = view.getUint16(4)
  const ancount = view.getUint16(6)

  let offset = 12 // Skip header

  // Skip question section
  for (let i = 0; i < qdcount; i++) {
    while (uint8[offset] !== 0) offset += uint8[offset] + 1
    offset += 5 // null label + QTYPE + QCLASS
  }

  // Parse answer section
  const answers: any[] = []

  for (let i = 0; i < ancount; i++) {
    // Name (may be compressed)
    if ((uint8[offset] & 0xc0) === 0xc0) {
      offset += 2 // Compressed name pointer
    } else {
      while (uint8[offset] !== 0) offset += uint8[offset] + 1
      offset += 1
    }

    const rtype = view.getUint16(offset); offset += 2
    const rclass = view.getUint16(offset); offset += 2
    const ttl = view.getUint32(offset); offset += 4
    const rdlength = view.getUint16(offset); offset += 2

    const rdata = uint8.slice(offset, offset + rdlength); offset += rdlength

    answers.push({
      name: '',
      type: rtype,
      TTL: ttl,
      data: parseRdata(rtype, rdata, uint8, buffer)
    })
  }

  return {
    Status: status,
    Answer: answers
  }
}

// 解析 RDATA
function parseRdata(type: number, rdata: Uint8Array, fullPacket: Uint8Array, buffer: ArrayBuffer): string {
  switch (type) {
    case 1: // A
      return Array.from(rdata).join('.')
    case 28: // AAAA
      const groups: string[] = []
      for (let i = 0; i < 16; i += 2) {
        groups.push(((rdata[i] << 8) | rdata[i + 1]).toString(16))
      }
      return groups.join(':')
    case 5: // CNAME
    case 2: // NS
      return parseDomainName(rdata, 0, fullPacket)
    case 15: // MX
      const preference = (rdata[0] << 8) | rdata[1]
      const exchange = parseDomainName(rdata, 2, fullPacket)
      return `${preference} ${exchange}`
    case 16: // TXT
      let txt = ''
      let i = 0
      while (i < rdata.length) {
        const len = rdata[i]
        txt += new TextDecoder().decode(rdata.slice(i + 1, i + 1 + len))
        i += len + 1
      }
      return txt
    case 6: // SOA
      let off = 0
      const mname = parseDomainName(rdata, off, fullPacket)
      off = skipDomainName(rdata, off)
      const rname = parseDomainName(rdata, off, fullPacket)
      return `${mname} ${rname}`
    default:
      return Array.from(rdata).map(b => b.toString(16).padStart(2, '0')).join(' ')
  }
}

// 解析域名（支持压缩指针）
function parseDomainName(data: Uint8Array, offset: number, fullPacket: Uint8Array): string {
  const labels: string[] = []
  let jumped = false
  let jumpedOffset = 0

  while (offset < data.length) {
    const len = data[offset]

    if (len === 0) break

    if ((len & 0xc0) === 0xc0) {
      // Compressed pointer
      if (!jumped) {
        jumped = true
        jumpedOffset = offset + 2
      }
      const ptr = ((len & 0x3f) << 8) | data[offset + 1]
      offset = ptr
      continue
    }

    offset++
    labels.push(new TextDecoder().decode(data.slice(offset, offset + len)))
    offset += len
  }

  return labels.join('.')
}

// 跳过域名
function skipDomainName(data: Uint8Array, offset: number): number {
  while (offset < data.length && data[offset] !== 0) {
    if ((data[offset] & 0xc0) === 0xc0) {
      return offset + 2
    }
    offset += data[offset] + 1
  }
  return offset + 1
}

// DNS 状态码转换
function getDnsStatus(status: number): string {
  const statusMap: Record<number, string> = {
    0: '成功',
    1: '格式错误',
    2: '服务器错误',
    3: '域名不存在',
    4: '未实现',
    5: '拒绝'
  }
  return statusMap[status] || `未知错误 (${status})`
}

// 记录类型颜色
function getRecordColor(type: string): string {
  const colors: Record<string, string> = {
    A: 'bg-blue-500/20 text-blue-400',
    AAAA: 'bg-green-500/20 text-green-400',
    CNAME: 'bg-purple-500/20 text-purple-400',
    MX: 'bg-yellow-500/20 text-yellow-400',
    TXT: 'bg-orange-500/20 text-orange-400',
    NS: 'bg-cyan-500/20 text-cyan-400',
    SOA: 'bg-pink-500/20 text-pink-400'
  }
  return colors[type] || 'bg-gray-500/20 text-gray-400'
}

// 复制值
async function copyValue(value: string) {
  await navigator.clipboard.writeText(value)
  appStore.showToast('success', '已复制')
}

// 加载历史
onMounted(() => {
  const saved = localStorage.getItem('dnsHistory')
  if (saved) {
    try {
      history.value = JSON.parse(saved)
    } catch {
      // 忽略
    }
  }
})
</script>
