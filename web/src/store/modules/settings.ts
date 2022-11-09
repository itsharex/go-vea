import { defineStore } from 'pinia'
import { SettingState } from '@/types/store/setting'
import defaultSettings from '@/settings'
import { localStorage } from '@/utils/storage'
import { useDynamicTitle } from '@/utils/dynamicTitle'

const { sideTheme, showSettings, topNav, tagsView, fixedHeader, sidebarLogo, dynamicTitle } = defaultSettings

const storageSetting = JSON.parse(localStorage.get('layout-setting')) || ''
const el = document.documentElement

export const useSettingStore = defineStore({
  id: 'setting',
  state: (): SettingState => ({
    title: '',
    theme: storageSetting.theme || getComputedStyle(el).getPropertyValue(`--el-color-primary`),
    sideTheme: storageSetting.sideTheme || sideTheme,
    showSettings: showSettings,
    topNav: storageSetting.topNav === undefined ? topNav : storageSetting.topNav,
    tagsView: storageSetting.tagsView === undefined ? tagsView : storageSetting.tagsView,
    fixedHeader: storageSetting.fixedHeader === undefined ? fixedHeader : storageSetting.fixedHeader,
    sidebarLogo: storageSetting.sidebarLogo === undefined ? sidebarLogo : storageSetting.sidebarLogo,
    dynamicTitle: storageSetting.dynamicTitle === undefined ? dynamicTitle : storageSetting.dynamicTitle
  }),
  actions: {
    async changeSetting(payload: { key: string; value: any }) {
      const { key, value } = payload
      switch (key) {
        case 'sideTheme':
          this.sideTheme = value
          break
        case 'topNav':
          this.topNav = value
          break
        case 'dynamicTitle':
          this.dynamicTitle = value
          break
        case 'theme':
          this.theme = value
          break
        case 'showSettings':
          this.showSettings = value
          break
        case 'fixedHeader':
          this.fixedHeader = value
          break
        case 'tagsView':
          this.tagsView = value
          localStorage.set('tagsView', value)
          break
        case 'sidebarLogo':
          this.sidebarLogo = value
          break
        default:
          break
      }
    },
    // 设置网页标题
    setTitle(title: string) {
      this.title = title
      useDynamicTitle()
    }
  }
})

export default useSettingStore
