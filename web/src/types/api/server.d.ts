export interface Cpu {
  cpuNum: string
  used: string
  free: string
  total: string
  sys: string
}
export interface Disk {
  dirName: string
  fsType: string
  total: string
  used: string
  free: string
  usedPercent: string
}
export interface Sys {
  computerName: string
  computerIp: string
  osName: string
  osArch: string
}
export interface Mem {
  total: string
  used: string
  free: string
  usedPercent: string
}

export type ServerFormData = {
  cpu: Cpu
  mem: Mem
  disk: Disk[]
  sys: Sys
}
