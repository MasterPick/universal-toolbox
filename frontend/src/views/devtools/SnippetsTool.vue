<template>
  <div class="page-container">
    <div class="page-title"><BookMarked :size="20" class="text-primary-400"/>代码片段管理</div>
    <div class="page-desc">保存和管理常用代码片段</div>
    <div class="flex gap-2 mb-4">
      <input v-model="searchKw" class="input-field flex-1" placeholder="搜索片段..." @input="loadList"/>
      <button @click="exportSnippets" class="btn btn-secondary"><Download :size="14"/>导出</button>
      <button @click="triggerImport" class="btn btn-secondary"><Upload :size="14"/>导入</button>
      <input ref="importInput" type="file" accept=".json" class="hidden" @change="importSnippets"/>
      <button @click="openAdd" class="btn btn-primary"><Plus :size="14"/>新增片段</button>
    </div>
    <div class="flex-1 overflow-auto space-y-2">
      <div v-for="s in snippets" :key="s.id" class="card group">
        <div class="flex items-start justify-between">
          <div>
            <div class="font-medium text-sm">{{ s.title }}</div>
            <div class="text-xs opacity-50 mt-0.5">{{ s.language }} · {{ s.createdAt }}</div>
          </div>
          <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button @click="copySnippet(s.content)" class="btn btn-secondary py-0.5 px-2 text-xs"><Copy :size="11"/></button>
            <button @click="editSnippet(s)" class="btn btn-secondary py-0.5 px-2 text-xs"><Edit3 :size="11"/></button>
            <button @click="deleteSnippet(s.id)" class="btn btn-danger py-0.5 px-2 text-xs"><Trash2 :size="11"/></button>
          </div>
        </div>
        <pre class="code-output mt-2 text-xs max-h-24 overflow-hidden">{{ s.content }}</pre>
      </div>
      <div v-if="snippets.length===0" class="text-center opacity-30 pt-16">
        <BookMarked :size="40" class="mx-auto mb-2"/>
        <div>暂无片段，点击"新增片段"添加</div>
      </div>
    </div>
    <!-- 新增/编辑弹窗 -->
    <div v-if="showEdit" class="modal-overlay" @click.self="showEdit=false">
      <div class="modal-box card w-[500px]">
        <div class="font-semibold mb-3">{{ editingId ? '编辑代码片段' : '新增代码片段' }}</div>
        <input v-model="newSnippet.title"    class="input-field mb-2" placeholder="片段标题..." />
        <input v-model="newSnippet.language" class="input-field mb-2" placeholder="编程语言（如 javascript）..." />
        <input v-model="newSnippet.tags"     class="input-field mb-2" placeholder="标签（逗号分隔）..." />
        <textarea v-model="newSnippet.content" class="textarea-field mb-3" rows="6" placeholder="代码内容..." />
        <div class="flex gap-2 justify-end">
          <button @click="showEdit=false" class="btn btn-secondary">取消</button>
          <button @click="saveSnippet" class="btn btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { BookMarked, Plus, Copy, Trash2, Edit3, Download, Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetSnippets, SaveSnippet, DeleteSnippet, SearchSnippets } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const snippets = ref<any[]>([])
const searchKw = ref('')
const showEdit  = ref(false)
const editingId = ref<number | null>(null)
const newSnippet = ref({ title:'', content:'', language:'text', tags:'' })
const importInput = ref<HTMLInputElement | null>(null)

async function loadList() {
  try {
    snippets.value = searchKw.value
      ? await SearchSnippets(searchKw.value) as any[]
      : await GetSnippets() as any[]
  } catch (e) {
    appStore.showToast('error', '加载片段失败: ' + String(e))
  }
}

function openAdd() {
  editingId.value = null
  newSnippet.value = { title:'', content:'', language:'text', tags:'' }
  showEdit.value = true
}

function editSnippet(s: any) {
  editingId.value = s.id
  newSnippet.value = { title: s.title, content: s.content, language: s.language, tags: s.tags }
  showEdit.value = true
}

async function saveSnippet() {
  const { title, content, language, tags } = newSnippet.value
  if (!title || !content) { appStore.showToast('warning', '标题和内容不能为空'); return }
  try {
    await SaveSnippet(editingId.value || 0, title, content, language, tags)
    showEdit.value = false
    newSnippet.value = { title:'', content:'', language:'text', tags:'' }
    editingId.value = null
    loadList()
    appStore.showToast('success', '片段已保存')
  } catch (e) {
    appStore.showToast('error', '保存失败: ' + String(e))
  }
}

async function deleteSnippet(id: number) {
  try {
    await DeleteSnippet(id)
    loadList()
    appStore.showToast('success', '已删除')
  } catch (e) {
    appStore.showToast('error', '删除失败: ' + String(e))
  }
}

async function copySnippet(content: string) {
  try {
    await navigator.clipboard.writeText(content)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

function triggerImport() {
  importInput.value?.click()
}

async function exportSnippets() {
  try {
    const data = JSON.stringify(snippets.value, null, 2)
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `snippets_${new Date().toISOString().slice(0,10)}.json`
    a.click()
    URL.revokeObjectURL(url)
    appStore.showToast('success', `已导出 ${snippets.value.length} 个片段`)
  } catch (e) {
    appStore.showToast('error', '导出失败: ' + String(e))
  }
}

async function importSnippets(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  try {
    const text = await file.text()
    const data = JSON.parse(text)
    const items = Array.isArray(data) ? data : [data]
    let count = 0
    for (const item of items) {
      if (item.title && item.content) {
        await SaveSnippet(0, item.title, item.content, item.language || 'text', item.tags || '')
        count++
      }
    }
    loadList()
    appStore.showToast('success', `已导入 ${count} 个片段`)
  } catch (err) {
    appStore.showToast('error', '导入失败: 文件格式错误')
  }
  // 重置 input
  if (importInput.value) importInput.value.value = ''
}

onMounted(loadList)
</script>
<style scoped>
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 100;
}
</style>
