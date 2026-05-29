import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api/v1": {
        target: "https://188.244.6.119:8888",
        changeOrigin: true,
        secure: false
      }
    }
  }
});
