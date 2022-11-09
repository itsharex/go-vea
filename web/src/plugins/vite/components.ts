import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import IconsResolver from 'unplugin-icons/resolver'

export default function createComponents() {
  return Components({
    // relative paths to the directory to search for components.
    dirs: ['src/components', 'src/views/components', 'src/layout/components'],

    // valid file extensions for components.
    extensions: ['vue'],
    // search for subdirectories
    deep: true,
    // resolvers for custom components 自动注册图标组件
    resolvers: [
      IconsResolver({
        prefix: false,
        enabledCollections: ['ep']
      }),
      ElementPlusResolver()
    ],

    // generate `components.d.ts` global declarations,
    // also accepts a path for custom filename
    dts: 'src/types/components.d.ts',

    // Allow subdirectories as namespace prefix for components.
    directoryAsNamespace: false,
    // Subdirectory paths for ignoring namespace prefixes
    // works when `directoryAsNamespace: true`
    globalNamespaces: [],

    // auto import for directives
    // default: `true` for Vue 3, `false` for Vue 2
    // Babel is needed to do the transformation for Vue 2, it's disabled by default for performance concerns.
    // To install Babel, run: `npm install -D @babel/parser @babel/traverse`
    directives: true,

    // filters for transforming targets
    include: [/\.vue$/, /\.vue\?vue/],
    exclude: [/[\\/]node_modules[\\/]/, /[\\/]\.git[\\/]/, /[\\/]\.nuxt[\\/]/]
  })
}
