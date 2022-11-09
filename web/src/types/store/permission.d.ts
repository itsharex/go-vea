/**
 * 权限类型声明
 */
export interface PermissionState {
  routes: RouteRecordRaw[]
  addRoutes: RouteRecordRaw[]
  defaultRoutes: RouteRecordRaw[]
  topbarRouters: RouteRecordRaw[]
  sidebarRouters: RouteRecordRaw[]
}
