import { login, logout, getUserInfo } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { LoginFormData } from '@/types/api/user'
import { UserState } from '@/types/store/user'
import defAva from '@/assets/images/avatar.svg'

const useUserStore = defineStore({
  id: 'user',
  state: (): UserState => ({
    token: getToken() || '',
    name: '',
    avatar: '',
    roles: [],
    perms: []
  }),
  actions: {
    async RESET_STATE() {
      this.$reset()
    },
    // 登录
    login(userInfo: LoginFormData) {
      const username = userInfo.username.trim()
      const password = userInfo.password
      const code = userInfo.code
      const uuid = userInfo.uuid
      return new Promise((resolve, reject) => {
        login({ username, password, code, uuid })
          .then(res => {
            const accessToken = res.msg
            setToken(accessToken)
            this.token = accessToken
            resolve(accessToken)
          })
          .catch(error => {
            reject(error)
          })
      })
    },
    // 获取用户信息
    getUserInfo() {
      return new Promise((resolve, reject) => {
        getUserInfo()
          .then(r => {
            const res = r.data
            const user = res.user
            const avatar = user.avatar == '' || user.avatar == null ? defAva : import.meta.env.VITE_APP_BASE_API + user.avatar

            if (res.roles && res.roles.length > 0) {
              // 验证返回的roles是否是一个非空数组
              this.roles = res.roles
              this.perms = res.permissions
            } else {
              this.roles = ['ROLE_DEFAULT']
            }
            this.name = user.userName
            this.avatar = avatar
            resolve(res)
          })
          .catch(error => {
            reject(error)
          })
      })
    },
    // 退出系统
    logout() {
      return new Promise((resolve, reject) => {
        logout()
          .then(() => {
            this.token = ''
            this.roles = []
            this.perms = []
            removeToken()
            resolve(null)
          })
          .catch(error => {
            reject(error)
          })
      })
    },

    // 清除 Token
    resetToken() {
      return new Promise(resolve => {
        removeToken()
        this.RESET_STATE()
        resolve(null)
      })
    }
  }
})

export default useUserStore
