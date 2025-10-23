import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [['babel-plugin-react-compiler']],
      },
    }),
    tailwindcss()
  ],
  server: {
    proxy: {
      '/api': {
        changeOrigin: true,
        target: 'http://localhost:8080/',
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      }
    }
  }
})
