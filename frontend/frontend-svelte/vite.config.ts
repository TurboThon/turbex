import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import wasmPack from 'vite-plugin-wasm-pack';

export default defineConfig({
	plugins: [sveltekit(), wasmPack(['../turbex-crypt'])],
  server: {
    proxy: {
      '/api': {
       target: 'http://localhost:8000',
       changeOrigin: true,
       secure: false,      
       ws: true,
      }
    }
  }
});
