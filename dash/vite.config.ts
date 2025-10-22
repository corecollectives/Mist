import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [['babel-plugin-react-compiler']],
      },
    }),
  ],
  server: {
    proxy: {
      '/api': {
        changeOrigin: true,
<<<<<<< HEAD
        target: process.env.API_URL || 'http://localhost:8080',
=======
        target: 'http://localhost:8080/',
>>>>>>> websockets
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      }
    }
  }
})
