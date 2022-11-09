import { defineConfig, loadEnv } from 'vite'
import createVitePlugins from './src/plugins/vite'
import { resolve } from 'path' // 编辑器提示 path 模块找不到，可以yarn add @types/node --dev

// https://vitejs.dev/config/
export default defineConfig(({ mode, command }) => {
  const env = loadEnv(mode, process.cwd())
  return {
    plugins: createVitePlugins(env, command === 'build'),
    resolve: {
      // https://cn.vitejs.dev/config/#resolve-alias
      alias: {
        // 设置路径
        '~': resolve(__dirname, './'),
        // 设置别名
        '@': resolve(__dirname, './src'),
        // 国际化
        'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js'
      },
      // https://cn.vitejs.dev/config/#resolve-extensions
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
    },
    server: {
      port: Number(env.VITE_APP_PORT), // 设置服务启动端口号
      open: true, // 设置服务启动时是否自动打开浏览器
      // 代理
      proxy: {
        [env.VITE_APP_BASE_API]: {
          target: 'http://localhost:8081',
          changeOrigin: true,
          rewrite: path => path.replace(new RegExp(`^${env.VITE_APP_BASE_API}`), '')
        }
      }
    }
  }
})
