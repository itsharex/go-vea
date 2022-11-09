/**
 * v-hasPerms 操作权限处理
 */
import useStore from '@/store'

export default {
  mounted(el: HTMLElement, binding: any, vnode: any) {
    const { value } = binding
    const all_permission = '*:*:*'
    const perms = useStore().user.perms

    if (value && value instanceof Array && value.length > 0) {
      const permissionFlag = value

      const hasPermissions = perms.some(p => {
        return all_permission === p || permissionFlag.includes(p)
      })

      if (!hasPermissions) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    } else {
      throw new Error(`请设置操作权限标签值`)
    }
  }
}
