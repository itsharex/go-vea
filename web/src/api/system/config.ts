import request from '@/utils/request'

const server_name = '/system/config'

// 查询参数列表
export function listConfig(data) {
  return request({
    url: `${server_name}/list`,
    method: 'post',
    data: data
  })
}

// 查询参数详细
export function getConfig(configId) {
  return request({
    url: `${server_name}/` + configId,
    method: 'get'
  })
}

// 根据参数键名查询参数值
export function getConfigKey(configKey) {
  return request({
    url: `${server_name}/configKey/` + configKey,
    method: 'get'
  })
}

// 新增参数配置
export function addConfig(data) {
  return request({
    url: `${server_name}`,
    method: 'post',
    data: data
  })
}

// 修改参数配置
export function updateConfig(data) {
  return request({
    url: `${server_name}`,
    method: 'put',
    data: data
  })
}

// 删除参数配置
export function delConfig(ids) {
  return request({
    url: `${server_name}`,
    method: 'delete',
    data: ids
  })
}

// 刷新参数缓存
export function refreshCache() {
  return request({
    url: `${server_name}/refreshCache`,
    method: 'delete'
  })
}
