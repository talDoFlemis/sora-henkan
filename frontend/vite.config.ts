import { defineConfig } from "vite"
import react from "@vitejs/plugin-react-swc"
import { viteEnvs } from "vite-envs"
import tailwindcss from "@tailwindcss/vite"
import path from "path"

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    viteEnvs({
      declarationFile: ".env.example",
    }),
    react(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
