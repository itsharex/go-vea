import request from '@/utils/request'

// 查询登录日志列表
export function list(data) {
  return request({
    url: '/monitor/loginLog/list',
    method: 'post',
    data: data
  })
}

// 删除登录日志
export function delLogininfor(infoId) {
  return request({
    url: '/monitor/loginLog/' + infoId,
    method: 'delete'
  })
}

// 解锁用户登录状态
export function unlockLogininfor(username) {
  return request({
    url: '/monitor/loginLog/unlock/' + username,
    method: 'get'
  })
}

// 清空登录日志
export function cleanLogininfor() {
  return request({
    url: '/monitor/loginLog/clean',
    method: 'delete'
  })
}
