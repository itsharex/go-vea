package request

import (
	"go-vea/app/common/base"
	"time"
)

type SysOperLog struct {
	base.CommonModel
	OperID        int64      `json:"operId"`        // 日志主键
	Title         string     `json:"title"`         // 模块标题
	BusinessType  int64      `json:"businessType"`  // 业务类型（0其它 1新增 2修改 3删除）
	BusinessTypes []int64    `json:"businessTypes"` // 业务类型（0其它 1新增 2修改 3删除）
	Method        string     `json:"method"`        // 方法名称
	RequestMethod string     `json:"requestMethod"` // 请求方式
	OperatorType  int64      `json:"operatorType"`  // 操作类别（0其它 1后台用户 2手机端用户）
	OperName      string     `json:"operName"`      // 操作人员
	DeptName      string     `json:"deptName"`      // 部门名称
	OperURL       string     `json:"operUrl"`       // 请求URL
	OperIP        string     `json:"operIp"`        // 主机地址
	OperLocation  string     `json:"operLocation"`  // 操作地点
	OperParam     string     `json:"operParam"`     // 请求参数
	JSONResult    string     `json:"jsonResult"`    // 返回参数
	Status        int64      `json:"status"`        // 操作状态（0正常 1异常）
	ErrorMsg      string     `json:"errorMsg"`      // 错误消息
	OperTime      *time.Time `json:"operTime"`      // 操作时间
}
