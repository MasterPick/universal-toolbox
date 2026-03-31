<template>
  <!-- Toast 通知容器 -->
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast" tag="div" class="flex flex-col gap-2">
        <div
          v-for="toast in appStore.toasts"
          :key="toast.id"
          :class="['toast-item', `toast-${toast.type}`]"
          @click="appStore.removeToast(toast.id)"
        >
          <!-- 图标 -->
          <component :is="iconMap[toast.type]" :size="16" class="shrink-0" />
          <!-- 消息文本 -->
          <span class="text-sm">{{ toast.message }}</span>
          <!-- 关闭按钮 -->
          <X :size="13" class="shrink-0 opacity-60 hover:opacity-100 cursor-pointer ml-auto" />
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { CheckCircle2, XCircle, AlertTriangle, Info, X } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 各类型对应图标
const iconMap = {
  success: CheckCircle2,
  error: XCircle,
  warning: AlertTriangle,
  info: Info,
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 9999;
  pointer-events: none;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  border-radius: 10px;
  min-width: 260px;
  max-width: 400px;
  pointer-events: all;
  cursor: pointer;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  color: #fff;
}

.toast-success { background: rgba(16, 185, 129, 0.85); }
.toast-error   { background: rgba(239, 68, 68, 0.85); }
.toast-warning { background: rgba(245, 158, 11, 0.85); }
.toast-info    { background: rgba(99, 102, 241, 0.85); }

/* 进入/离开动画 */
.toast-enter-active { transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1); }
.toast-leave-active { transition: all 0.2s ease-in; }
.toast-enter-from   { opacity: 0; transform: translateY(16px) scale(0.95); }
.toast-leave-to     { opacity: 0; transform: translateX(20px); }
</style>
