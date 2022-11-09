import request from '@/utils/request'
import { LoginFormData } from '@/types/api/user'

// 登录方法
export function login(data: LoginFormData) {
  return request({
    url: '/login',
    headers: {
      isToken: false
    },
    method: 'post',
    data: data
  })
}

// 注册方法
export function register(data) {
  return request({
    url: '/register',
    headers: {
      isToken: false
    },
    method: 'post',
    data: data
  })
}

// 获取用户详细信息
export function getUserInfo() {
  return request({
    url: '/getUserInfo',
    method: 'get'
  })
}

// 退出方法
export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}

// 获取验证码
export function getCodeImg() {
  return request({
    url: '/captchaImage',
    headers: {
      isToken: false
    },
    method: 'get',
    timeout: 20000
  })
}

// 获取路由
export function getRouters() {
  return request({
    url: '/getRouters',
    method: 'get'
  })
}
