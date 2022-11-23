import request from '@/utils/request'

// 查询操作日志列表
export function list(data) {
  return request({
    url: '/monitor/operLog/list',
    method: 'post',
    data: data
  })
}

// 删除操作日志
export function delOperlog(data) {
  return request({
    url: '/monitor/operLog',
    method: 'delete',
    data: data
  })
}

// 清空操作日志
export function cleanOperlog() {
  return request({
    url: '/monitor/operLog/clean',
    method: 'delete'
  })
}
