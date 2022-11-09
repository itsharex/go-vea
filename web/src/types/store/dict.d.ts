/**
 * 字典
 */
interface kv {
  value: string
  key: string
}

export interface DictState {
  dict: kv[]
}
