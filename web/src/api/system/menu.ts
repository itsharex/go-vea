import request from '@/utils/request'

// 查询菜单列表
export function listMenu(data) {
  return request({
    url: '/system/menu/list',
    method: 'post',
    data: data
  })
}

// 分页查询菜单列表
export function listTreeByPage(data) {
  return request({
    url: '/system/menu/listTreeByPage',
    method: 'post',
    params: data
  })
}

// 查询菜单详细
export function getMenu(menuId) {
  return request({
    url: '/system/menu/' + menuId,
    method: 'get'
  })
}

// 查询菜单下拉树结构
export function treeSelect(data) {
  return request({
    url: '/system/menu/treeSelect',
    method: 'post',
    data: data
  })
}

// 根据角色ID查询菜单下拉树结构
export function roleMenuTreeselect(data) {
  return request({
    url: '/system/menu/roleMenuTreeSelect',
    method: 'post',
    data: data
  })
}

// 新增菜单
export function addMenu(data) {
  return request({
    url: '/system/menu',
    method: 'post',
    data: data
  })
}

// 修改菜单
export function updateMenu(data) {
  return request({
    url: '/system/menu',
    method: 'put',
    data: data
  })
}

// 删除菜单
export function delMenu(menuId) {
  return request({
    url: '/system/menu/' + menuId,
    method: 'delete'
  })
}
