import { viteMockServe } from 'vite-plugin-mock'

export default function createMockServer() {
  return viteMockServe({
    // 解析根目录下的mock文件夹
    mockPath: 'mock',
    // 开发打包开关
    // localEnabled: localEnabled,
    // 生产打包开关
    // prodEnabled: prodEnabled,
    // 打开后，可以读取 ts 文件模块。 请注意，打开后将无法监视.js 文件。
    supportTs: true,
    // 监视文件更改
    watchFiles: true
  })
}
