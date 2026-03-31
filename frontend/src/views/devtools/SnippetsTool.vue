<template>
  <div class="page-container">
    <div class="page-title"><BookMarked :size="20" class="text-primary-400"/>代码片段管理</div>
    <div class="page-desc">保存和管理常用代码片段</div>
    <div class="flex gap-2 mb-4">
      <input v-model="searchKw" class="input-field flex-1" placeholder="搜索片段..." @input="loadList"/>
      <button @click="showAdd=true" class="btn btn-primary"><Plus :size="14"/>新增片段</button>
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
    <!-- 新增弹窗 -->
    <div v-if="showAdd" class="modal-overlay" @click.self="showAdd=false">
      <div class="modal-box card w-[500px]">
        <div class="font-semibold mb-3">新增代码片段</div>
        <input v-model="newSnippet.title"    class="input-field mb-2" placeholder="片段标题..." />
        <input v-model="newSnippet.language" class="input-field mb-2" placeholder="编程语言（如 javascript）..." />
        <input v-model="newSnippet.tags"     class="input-field mb-2" placeholder="标签（逗号分隔）..." />
        <textarea v-model="newSnippet.content" class="textarea-field mb-3" rows="6" placeholder="代码内容..." />
        <div class="flex gap-2 justify-end">
          <button @click="showAdd=false" class="btn btn-secondary">取消</button>
          <button @click="saveSnippet" class="btn btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { BookMarked, Plus, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetSnippets, SaveSnippet, DeleteSnippet, SearchSnippets } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const snippets = ref<any[]>([])
const searchKw = ref('')
const showAdd  = ref(false)
const newSnippet = ref({ title:'', content:'', language:'text', tags:'' })

async function loadList() {
  snippets.value = searchKw.value
    ? await SearchSnippets(searchKw.value) as any[]
    : await GetSnippets() as any[]
}

async function saveSnippet() {
  const { title, content, language, tags } = newSnippet.value
  if (!title || !content) { appStore.showToast('warning', '标题和内容不能为空'); return }
  await SaveSnippet(title, content, language, tags)
  showAdd.value = false
  newSnippet.value = { title:'', content:'', language:'text', tags:'' }
  loadList()
  appStore.showToast('success', '片段已保存')
}

async function deleteSnippet(id: number) {
  await DeleteSnippet(id)
  loadList()
  appStore.showToast('success', '已删除')
}

async function copySnippet(content: string) {
  await navigator.clipboard.writeText(content)
  appStore.showToast('success', '已复制到剪贴板')
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
