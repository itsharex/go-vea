import request from '@/utils/request'

const server_name = '/system/dept'

// 查询部门列表
export function listDept(data) {
  return request({
    url: `${server_name}/list`,
    method: 'post',
    data: data
  })
}

// 查询部门列表
export function listDeptTree(data) {
  return request({
    url: `${server_name}/tree`,
    method: 'post',
    params: data
  })
}

// 查询部门列表（排除节点）
export function listDeptExcludeChild(deptId) {
  return request({
    url: '/system/dept/list/exclude/' + deptId,
    method: 'get'
  })
}

// 查询部门详细
export function getDept(deptId) {
  return request({
    url: '/system/dept/' + deptId,
    method: 'get'
  })
}

// 新增部门
export function addDept(data) {
  return request({
    url: '/system/dept',
    method: 'post',
    data: data
  })
}

// 修改部门
export function updateDept(data) {
  return request({
    url: '/system/dept',
    method: 'put',
    data: data
  })
}

// 删除部门
export function delDept(deptId) {
  return request({
    url: '/system/dept/' + deptId,
    method: 'delete'
  })
}
