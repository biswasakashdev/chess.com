import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig, loadEnv } from "vite"

// https://vite.dev/config/

export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  // Set the third parameter to '' to load all env variables regardless of the 'VITE_' prefix.
  const env = loadEnv(mode, process.cwd(), "")
  const serverUrl = env.VITE_SERVER_URL

  return {
    server: {
      proxy: {
        "/backend/ws": {
          target: `ws://${serverUrl}`,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/backend/, ""),
        },
        "/backend/api": {
          target: `http://${serverUrl}`,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/backend/, ""),
        },
      },
    },
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(import.meta.dirname, "./src"),
      },
    },
  }
})
