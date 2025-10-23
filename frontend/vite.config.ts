import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { viteEnvs } from 'vite-envs'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    viteEnvs({
      declarationFile: '.env.example',
    }),
    react(),
  ],
})
