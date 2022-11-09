import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index'
import '@/permission' // 路由守卫

// svg图标
import 'virtual:svg-icons-register'

// 自定义样式
import '@/assets/styles/index.scss'

// 注册指令
import plugins from './plugins' // plugins

import { useDict } from '@/utils/dict'
import { parseTime, resetForm, addDateRange, handleTree, selectDictLabel, selectDictLabels } from '@/utils/ruoyi'

// 国际化
import i18n from './i18n/index'

import directive from './directive' // directive

const app = createApp(App)

// 全局方法挂载
app.config.globalProperties.useDict = useDict
app.config.globalProperties.parseTime = parseTime
app.config.globalProperties.resetForm = resetForm
app.config.globalProperties.handleTree = handleTree
app.config.globalProperties.addDateRange = addDateRange
app.config.globalProperties.selectDictLabel = selectDictLabel
app.config.globalProperties.selectDictLabels = selectDictLabels

// 创建pinia 实例
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(plugins)
app.use(i18n)

directive(app)

app.mount('#app')
