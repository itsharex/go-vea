import { defineConfig, loadEnv } from 'vite'
import createVitePlugins from './src/plugins/vite'
import { resolve } from 'path' // 编辑器提示 path 模块找不到，可以yarn add @types/node --dev

// https://vitejs.dev/config/
export default defineConfig(({ mode, command }) => {
  const env = loadEnv(mode, process.cwd())
  const { VITE_APP_ENV } = env
  return {
    // 按实际情况修改 打包路径（就是网站前缀，例如部署到 https://aguoxing.github.io/go-vea/ 域名下，就需要填写 /go-vea/）
    base: VITE_APP_ENV === 'production' ? '/' : '/',
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
    },
    build: {
      sourcemap: false,
      // 消除打包大小超过500kb警告
      chunkSizeWarningLimit: 4000,
      // Vite 2.6.x 以上需要配置 minify: "terser", terserOptions 才能生效
      // minify: "terser",
      // 在打包代码时移除 console.log、debugger 和 注释
      /*terserOptions: {
        compress: {
          drop_console: false,
          drop_debugger: true,
          pure_funcs: ["console.log"]
        },
        format: {
          // 删除注释
          comments: false
        }
      },*/
      rollupOptions: {
        input: {
          index: resolve("index.html")
        },
        // 静态资源分类打包
        output: {
          chunkFileNames: "static/js/[name]-[hash].js",
          entryFileNames: "static/js/[name]-[hash].js",
          assetFileNames: "static/[ext]/[name]-[hash].[ext]"
        }
      }
    }
  }
})
