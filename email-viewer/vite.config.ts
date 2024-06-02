import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'

const env = loadEnv('', process.cwd());

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VueDevTools(),
  ],
//   define: {
//     VITE_API_URL: env.VITE_API_URL,
//     VITE_API_USER: env.VITE_API_URL,
//     VITE_API_USER_PWD: env.VITE_API_USER_PWD
//   },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
