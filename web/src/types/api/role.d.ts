export type RoleFormData = {
  roleId: number | undefined
  roleName: string
  status: number | string
  roleKey: string
  roleSort: number
  menuCheckStrictly: boolean
  deptCheckStrictly: boolean
  remark: string
  dataScope: number
  deptIds: number[]
  menuIds: number[]
}
