export interface PageQueryParam {
  pageNum: number
  pageSize: number
}

export interface PageResult<T> {
  list: T
  total: number
}

export interface Ids {
  ids: string[] | number[]
}
