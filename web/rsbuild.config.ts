import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [pluginReact()],
  source: {
    entry: {
      index: './src/main.tsx',
    },
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  output: {
    distPath: {
      root: path.resolve(__dirname, '../resources/static'),
      js: 'js',
      css: 'css',
      svg: 'assets',
      image: 'assets',
      font: 'assets',
    },
    assetPrefix: '/',
  },
  server: {
    port: 3001,
    proxy: {
      '/api': {
        target: 'http://localhost:18080',
        changeOrigin: true,
      },
    },
  },
  html: {
    template: './index.html',
  },
});
