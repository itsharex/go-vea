import { PageQueryParam, PageResult } from './base'

/**
 * 登录表单
 */
export interface LoginFormData {
  username: string
  password: string
  uuid: string
  code: any
}

/**
 * 登录响应
 */
export interface LoginResponseData {
  access_token: string
  token_type: string
}

/**
 * 登录用户信息
 */
export interface UserInfo {
  nickname: string
  avatar: string
  roles: string[]
  perms: string[]
}

/**
 * 用户查询参数
 */
export interface UserQueryParam extends PageQueryParam {
  keywords: string
  status: number
  deptId: number
}

/**
 * 用户分页列表项声明
 */
export interface UserItem {
  id: string
  username: string
  nickname: string
  mobile: string
  gender: number
  avatar: string
  email: string
  status: number
  deptName: string
  roleNames: string
  createTime: string
}

/**
 * 用户分页项类型声明
 */
export type UserPageResult = PageResult<UserItem[]>

/**
 * 用户表单类型声明
 */
export interface UserFormData {
  userId: number | undefined
  postId: number | undefined
  deptId: number | undefined
  username: string | undefined
  nickname: string | undefined
  gender: string | undefined
  password: string | undefined
  phoneNumber: string | undefined
  email: string | undefined
  gender: number | undefined
  status: string | undefined
  remark: string | undefined
  postIds: number[] | undefined
  roleIds: number[] | undefined
}

/**
 * 用户导入表单类型声明
 */
export interface UserImportData {
  deptId: number
  roleIds: number[]
}
