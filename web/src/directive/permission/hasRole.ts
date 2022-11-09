/**
 * v-hasRole 角色权限处理
 */
import useStore from '@/store'

export default {
  mounted(el: HTMLElement, binding: any, vnode: any) {
    const { value } = binding
    const super_admin = 'admin'
    const { user } = useStore()

    if (value && value instanceof Array && value.length > 0) {
      const roleFlag = value

      const hasRole = user.roles.some(role => {
        return super_admin === role || roleFlag.includes(role)
      })

      if (!hasRole) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    } else {
      throw new Error(`请设置角色权限标签值`)
    }
  }
}
