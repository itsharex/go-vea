import useUserStore from './modules/user'
import useDictStore from './modules/dict'
import useAppStore from './modules/app'
import usePermissionStore from './modules/permission'
import useSettingStore from './modules/settings'
import useTagsViewStore from './modules/tagsView'

export const store = createPinia()

const useStore = () => ({
  user: useUserStore(),
  dict: useDictStore(),
  app: useAppStore(),
  permission: usePermissionStore(),
  setting: useSettingStore(),
  tagsView: useTagsViewStore()
})

export default useStore
