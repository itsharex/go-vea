export type MenuFormData = {
  menuId: number
  parentId: number
  menuType: string
  menuName: string
  orderNum: number
  isFrame: string | number
  path: string
  component: string
  perms: string
  query: string
  icon: string
  isCache: string | number
  visible: string
  status: string
  children: []
}

export type MenuOptions = {
  menuId: number
  menuName: string
  children: any[]
}
