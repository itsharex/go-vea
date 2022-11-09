import router from './router'
import useStore from './store'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { getToken } from './utils/auth'
NProgress.configure({ showSpinner: false }) // 进度环显示/隐藏

// 白名单路由
const whiteList = ['/login']

router.beforeEach(async (to, from, next) => {
  NProgress.start()
  const { user, permission, setting } = useStore()
  // const hasToken = user.token
  const hasToken = getToken()
  if (hasToken) {
    to.meta.title && setting.setTitle(to.meta.title)
    // 登录成功，跳转到首页
    if (to.path === '/login') {
      next({ path: '/' })
      NProgress.done()
    } else {
      const hasGetUserInfo = user.roles.length > 0
      if (hasGetUserInfo) {
        if (to.matched.length === 0) {
          from.name ? next({ name: from.name as any }) : next('/401')
        } else {
          next()
        }
      } else {
        try {
          await user.getUserInfo()
          const roles = user.roles
          // 根据roles权限生成可访问的路由表
          const accessRoutes: any = await permission.generateRoutes(roles)
          accessRoutes.forEach((route: any) => {
            router.addRoute(route) // 动态添加可访问路由表
          })
          next({ ...to, replace: true }) // hack方法 确保addRoutes已完成
        } catch (error) {
          // 移除 token 并跳转登录页
          await user.resetToken()
          ElMessage.error((error as any) || 'Has Error')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    // 未登录可以访问白名单页面(登录页面)
    if (whiteList.indexOf(to.path) !== -1) {
      next()
    } else {
      next(`/login?redirect=${to.fullPath}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
