package response

type Sys struct {
	ComputerName string `json:"computerName"` // 服务器名称
	ComputerIp   string `json:"computerIp"`   // 服务器Ip
	UserDir      string `json:"userDir"`      // 项目路径
	OsName       string `json:"osName"`       // 操作系统
	OsArch       string `json:"osArch"`       // 系统架构
}

type Disk struct {
	DirName     string `json:"dirName"`     // 路径
	FsType      string `json:"fsType"`      // 文件类型
	Total       string `json:"total"`       // 总大小
	Used        string `json:"used"`        // 已用大小
	Free        string `json:"free"`        // 可用大小
	UsedPercent string `json:"usedPercent"` // 已使用百分比
}

type Cpu struct {
	CpuNum int32   `json:"cpuNum"` // 核心数
	Total  float64 `json:"total"`  // CPU总的使用率
	Sys    float64 `json:"sys"`    // CPU系统使用率
	Used   float64 `json:"used"`   // CPU用户使用率
	Wait   float64 `json:"wait"`   // CPU当前等待率
	Free   float64 `json:"free"`   //CPU当前空闲率
}

type Mem struct {
	Total       uint64  `json:"total"`       // 内存总量
	Used        uint64  `json:"used"`        // 已用内存
	Free        uint64  `json:"free"`        // 剩余内存
	UsedPercent float64 `json:"usedPercent"` // 已使用百分比
}

type ServerInfo struct {
	Cpu  *Cpu    `json:"cpu"`
	Sys  *Sys    `json:"sys"`
	Mem  *Mem    `json:"mem"`
	Disk []*Disk `json:"disk"`
}
