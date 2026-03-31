<template>
  <div class="page-container">
    <div class="page-title"><StickyNote :size="20" class="text-primary-400"/>备忘录</div>
    <div class="page-desc">快速记录想法与待办事项</div>
    <div class="flex gap-2 mb-4">
      <button @click="openAdd" class="btn btn-primary"><Plus :size="14"/>新建备忘</button>
    </div>
    <div class="flex-1 overflow-auto grid grid-cols-3 gap-3 content-start">
      <div v-for="n in notes" :key="n.id" class="note-card group" :style="{ borderColor: n.color + '44' }">
        <div class="flex items-start justify-between mb-1">
          <div class="font-medium text-sm flex items-center gap-1">
            <Pin :size="12" v-if="n.pinned" class="text-yellow-400"/>
            {{ n.title }}
          </div>
          <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button @click="pinNote(n.id)" class="btn btn-secondary py-0.5 px-1.5"><Pin :size="11"/></button>
            <button @click="deleteNote(n.id)" class="btn btn-danger py-0.5 px-1.5"><Trash2 :size="11"/></button>
          </div>
        </div>
        <div class="text-xs opacity-70 leading-relaxed whitespace-pre-wrap line-clamp-4">{{ n.content }}</div>
        <div class="text-xs opacity-30 mt-2">{{ n.updatedAt }}</div>
      </div>
      <div v-if="notes.length===0" class="col-span-3 text-center opacity-30 pt-16">
        <StickyNote :size="40" class="mx-auto mb-2"/><div>暂无备忘，点击新建</div>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showEdit" class="modal-overlay" @click.self="showEdit=false">
      <div class="modal-box card w-[440px]">
        <div class="font-semibold mb-3">{{ editNote.id ? '编辑备忘' : '新建备忘' }}</div>
        <input v-model="editNote.title" class="input-field mb-2" placeholder="标题..." />
        <textarea v-model="editNote.content" class="textarea-field mb-3" rows="6" placeholder="内容..." />
        <div class="flex gap-2 justify-end">
          <button @click="showEdit=false" class="btn btn-secondary">取消</button>
          <button @click="saveNote" class="btn btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { StickyNote, Plus, Pin, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetNotes, SaveNote, PinNote, DeleteNote } from '../../../wailsjs/go/daily/DailyTools'
const appStore = useAppStore()
const notes = ref<any[]>([])
const showEdit = ref(false)
const editNote = ref({ id: 0, title: '', content: '', color: '#6366f1' })

async function loadNotes() { notes.value = await GetNotes() as any[] }
function openAdd() { editNote.value = { id:0, title:'', content:'', color:'#6366f1' }; showEdit.value = true }
async function saveNote() {
  await SaveNote(editNote.value.title, editNote.value.content, editNote.value.color)
  showEdit.value = false; loadNotes()
  appStore.showToast('success', '备忘已保存')
}
async function pinNote(id: number) { await PinNote(id); loadNotes() }
async function deleteNote(id: number) { await DeleteNote(id); loadNotes() }
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
