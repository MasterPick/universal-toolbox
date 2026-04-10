<template>
  <div class="page-container">
    <div class="page-title"><StickyNote :size="20" class="text-primary-400"/>备忘录</div>
    <div class="page-desc">快速记录想法与待办事项</div>
    <div class="flex gap-2 mb-4">
      <input v-model="searchKw" class="input-field flex-1" placeholder="搜索备忘录..."/>
      <button @click="openAdd" class="btn btn-primary"><Plus :size="14"/>新建备忘</button>
    </div>
    <div class="flex-1 overflow-auto grid grid-cols-3 gap-3 content-start">
      <div v-for="n in filteredNotes" :key="n.id" class="note-card group" :style="{ borderColor: n.color + '44' }">
        <div class="flex items-start justify-between mb-1">
          <div class="font-medium text-sm flex items-center gap-1">
            <Pin :size="12" v-if="n.pinned" class="text-yellow-400"/>
            {{ n.title }}
          </div>
          <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button @click="editNote(n)" class="btn btn-secondary py-0.5 px-1.5"><Edit3 :size="11"/></button>
            <button @click="pinNote(n.id)" class="btn btn-secondary py-0.5 px-1.5"><Pin :size="11"/></button>
            <button @click="deleteNote(n.id)" class="btn btn-danger py-0.5 px-1.5"><Trash2 :size="11"/></button>
          </div>
        </div>
        <div class="text-xs opacity-70 leading-relaxed whitespace-pre-wrap line-clamp-4">{{ n.content }}</div>
        <div class="text-xs opacity-30 mt-2">{{ n.updatedAt }}</div>
      </div>
      <div v-if="filteredNotes.length===0" class="col-span-3 text-center opacity-30 pt-16">
        <StickyNote :size="40" class="mx-auto mb-2"/><div>{{ searchKw ? '未找到匹配的备忘录' : '暂无备忘，点击新建' }}</div>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showEdit" class="modal-overlay" @click.self="showEdit=false">
      <div class="modal-box card w-[440px]">
        <div class="font-semibold mb-3">{{ editNoteData.id ? '编辑备忘' : '新建备忘' }}</div>
        <input v-model="editNoteData.title" class="input-field mb-2" placeholder="标题..." />
        <textarea v-model="editNoteData.content" class="textarea-field mb-3" rows="6" placeholder="内容..." />
        <div class="mb-3">
          <div class="label mb-2">颜色</div>
          <div class="flex gap-2">
            <button v-for="c in colorOptions" :key="c" @click="editNoteData.color = c"
              :class="['w-8 h-8 rounded-full border-2 transition-all', editNoteData.color === c ? 'border-white scale-110' : 'border-transparent']"
              :style="{ backgroundColor: c }"/>
          </div>
        </div>
        <div class="flex gap-2 justify-end">
          <button @click="showEdit=false" class="btn btn-secondary">取消</button>
          <button @click="saveNote" class="btn btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { StickyNote, Plus, Pin, Trash2, Edit3 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetNotes, SaveNote, UpdateNote, PinNote, DeleteNote } from '../../../wailsjs/go/daily/DailyTools'
const appStore = useAppStore()
const notes = ref<any[]>([])
const showEdit = ref(false)
const searchKw = ref('')
const editNoteData = ref({ id: 0, title: '', content: '', color: '#6366f1' })
const colorOptions = ['#6366f1', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ef4444', '#06b6d4']

const filteredNotes = computed(() => {
  if (!searchKw.value) return notes.value
  const kw = searchKw.value.toLowerCase()
  return notes.value.filter(n =>
    n.title.toLowerCase().includes(kw) || n.content.toLowerCase().includes(kw)
  )
})

async function loadNotes() {
  try {
    notes.value = await GetNotes() as any[]
  } catch (e) {
    appStore.showToast('error', '加载备忘录失败: ' + String(e))
  }
}

function openAdd() {
  editNoteData.value = { id: 0, title: '', content: '', color: '#6366f1' }
  showEdit.value = true
}

function editNote(n: any) {
  editNoteData.value = { id: n.id, title: n.title, content: n.content, color: n.color || '#6366f1' }
  showEdit.value = true
}

async function saveNote() {
  try {
    if (editNoteData.value.id) {
      await UpdateNote(editNoteData.value.id, editNoteData.value.title, editNoteData.value.content, editNoteData.value.color)
    } else {
      await SaveNote(editNoteData.value.title, editNoteData.value.content, editNoteData.value.color, false)
    }
    showEdit.value = false
    loadNotes()
    appStore.showToast('success', '备忘已保存')
  } catch (e) {
    appStore.showToast('error', '保存失败: ' + String(e))
  }
}

async function pinNote(id: number) {
  try {
    await PinNote(id)
    loadNotes()
  } catch (e) {
    appStore.showToast('error', '操作失败: ' + String(e))
  }
}

async function deleteNote(id: number) {
  try {
    await DeleteNote(id)
    loadNotes()
    appStore.showToast('success', '已删除')
  } catch (e) {
    appStore.showToast('error', '删除失败: ' + String(e))
  }
}

onMounted(loadNotes)
</script>
<style scoped>
.note-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px; padding: 12px;
  transition: all 0.2s;
}
.note-card:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.15); }
.modal-overlay { position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100; }
.line-clamp-4 { display:-webkit-box;-webkit-line-clamp:4;-webkit-box-orient:vertical;overflow:hidden; }
</style>
