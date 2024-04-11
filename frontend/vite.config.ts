import { sentryVitePlugin } from "@sentry/vite-plugin";
import path from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), sentryVitePlugin({
    org: "hs-flensburg-gitlab-classroom",
    project: "frontend"
  })],

  server: {
    host: true,
    proxy: {
      "/api": `http://${process.env.docker === "true" ? "backend" : "127.0.0.1"}:3000`,
    },
  },

  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },

  build: {
    sourcemap: true
  }
});