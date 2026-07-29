import { defineConfig } from "vite";
import glsl from "vite-plugin-glsl";

export default defineConfig({
    server: {
        proxy: {
            "/ws": {
                target: "ws://localhost:8080",
                ws: true
            }
        }
    },
    plugins: [glsl()]
});