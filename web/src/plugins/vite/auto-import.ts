import autoImport from 'unplugin-auto-import/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default function createAutoImport() {
  return autoImport({
    // 全局引入插件
    imports: [
      // presets
      'vue',
      'vue-router',
      'pinia',
      // custom
      {
        '@vueuse/core': [
          // named imports
          'useMouse', // import { useMouse } from '@vueuse/core',
          'useWindowSize',
          // alias
          ['useFetch', 'useMyFetch'] // import { useFetch as useMyFetch } from '@vueuse/core',
        ],
        axios: [
          // default imports
          ['default', 'axios'] // import { default as axios } from 'axios',
        ]
        /* '[package-name]': [
          '[import-names]',
          // alias
          ['[from]', '[alias]'],
        ], */
      }
    ],
    // eslint报错解决
    eslintrc: {
      enabled: true, // Default `false` 生成一次就可以，避免每次工程启动都生成
      filepath: './.eslintrc-auto-import.json', // Default `./.eslintrc-auto-import.json`
      globalsPropValue: true // Default `true`, (true | false | 'readonly' | 'readable' | 'writable' | 'writeable')
    },
    dts: 'src/types/auto-import.d.ts', // 生成 `auto-import.d.ts` 全局声明
    resolvers: [ElementPlusResolver()]
  })
}
