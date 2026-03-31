// 前端应用入口文件
// 初始化 Vue3 应用、路由、状态管理

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// 引入全局样式
import './assets/styles/main.css'

// 创建 Vue 应用实例
const app = createApp(App)

// 注册 Pinia 状态管理
app.use(createPinia())

// 注册 Vue Router
app.use(router)

// 挂载应用到 DOM
app.mount('#app')
