import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// Vite 构建配置
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      // 路径别名，方便导入
      '@': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@views': resolve(__dirname, 'src/views'),
      '@stores': resolve(__dirname, 'src/stores'),
    },
  },
  build: {
    // 生产构建输出目录
    outDir: 'dist',
    // 代码分割优化
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'codemirror': ['vue-codemirror', '@codemirror/lang-json'],
          'icons': ['lucide-vue-next'],
        },
      },
    },
    // 压缩优化
    minify: 'esbuild',
    target: 'esnext',
  },
})
