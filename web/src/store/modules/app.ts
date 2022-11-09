import { AppState } from '@/types/store/app'
import { getLanguage } from '@/i18n/index'
import { localStorage } from '@/utils/storage'

const useAppStore = defineStore({
  id: 'app',
  state: (): AppState => ({
    device: 'desktop',
    sidebar: {
      opened: localStorage.get('sidebarStatus') ? !!+localStorage.get('sidebarStatus') : true,
      withoutAnimation: false,
      hide: false
    },
    language: getLanguage(),
    size: localStorage.get('size') || 'default'
  }),
  actions: {
    toggleSidebar() {
      this.sidebar.opened = !this.sidebar.opened
      this.sidebar.withoutAnimation = false
      if (this.sidebar.opened) {
        localStorage.set('sidebarStatus', 1)
      } else {
        localStorage.set('sidebarStatus', 0)
      }
    },
    closeSideBar(withoutAnimation: any) {
      localStorage.set('sidebarStatus', 0)
      this.sidebar.opened = false
      this.sidebar.withoutAnimation = withoutAnimation
    },
    toggleDevice(device: string) {
      this.device = device
    },
    setSize(size: string) {
      this.size = size
      localStorage.set('size', size)
    },
    toggleSideBarHide(status: boolean) {
      this.sidebar.hide = status
    },
    setLanguage(language: string) {
      this.language = language
      localStorage.set('language', language)
    }
  }
})

export default useAppStore
