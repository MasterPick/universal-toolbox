/** @type {import('tailwindcss').Config} */
export default {
  // 扫描所有 Vue 和 TS 文件中的类名
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  // 深色模式：通过 class 控制（配合主题切换）
  darkMode: 'class',
  theme: {
    extend: {
      // 自定义颜色系统
      colors: {
        // 主色调
        primary: {
          50:  '#eef2ff',
          100: '#e0e7ff',
          200: '#c7d2fe',
          300: '#a5b4fc',
          400: '#818cf8',
          500: '#6366f1',
          600: '#4f46e5',
          700: '#4338ca',
          800: '#3730a3',
          900: '#312e81',
        },
        // 背景色（深色主题）
        dark: {
          bg:     '#0f0f17',
          card:   '#1a1a2e',
          border: '#2d2d44',
          hover:  '#252538',
          text:   '#e2e8f0',
          muted:  '#94a3b8',
        },
        // 背景色（浅色主题）
        light: {
          bg:     '#f8fafc',
          card:   '#ffffff',
          border: '#e2e8f0',
          hover:  '#f1f5f9',
          text:   '#1e293b',
          muted:  '#64748b',
        },
      },
      // 自定义动画
      animation: {
        'fade-in':     'fadeIn 0.2s ease-out',
        'slide-in':    'slideIn 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
        'slide-up':    'slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
        'scale-in':    'scaleIn 0.2s cubic-bezier(0.16, 1, 0.3, 1)',
        'shimmer':     'shimmer 1.5s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%':   { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideIn: {
          '0%':   { transform: 'translateX(-20px)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        },
        slideUp: {
          '0%':   { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        scaleIn: {
          '0%':   { transform: 'scale(0.95)', opacity: '0' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        shimmer: {
          '0%':   { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
      },
      // 自定义字体
      fontFamily: {
        sans: ['"Microsoft YaHei"', '"PingFang SC"', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', '"Cascadia Code"', '"Fira Code"', 'monospace'],
      },
      // 自定义圆角
      borderRadius: {
        'xl':  '0.75rem',
        '2xl': '1rem',
        '3xl': '1.5rem',
      },
      // 自定义阴影
      boxShadow: {
        'glass':   '0 8px 32px 0 rgba(31, 38, 135, 0.15)',
        'card':    '0 4px 20px rgba(0, 0, 0, 0.08)',
        'hover':   '0 8px 30px rgba(0, 0, 0, 0.12)',
        'glow':    '0 0 20px rgba(99, 102, 241, 0.3)',
      },
      // 背景模糊
      backdropBlur: {
        xs: '2px',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
