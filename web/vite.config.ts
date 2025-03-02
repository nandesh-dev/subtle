import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        extensions: ['.tsx', '.ts'],
        alias: {
            '@/src': path.resolve(__dirname, 'src'),
            '@/gen': path.resolve(__dirname, 'gen'),
        },
    },
    server: {
        port: 2001,
    },
    plugins: [react(), tailwindcss()],
})
