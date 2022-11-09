package e

var MsgFlags = map[int]string{
	SUCCESS:          "操作成功",
	ERROR:            "系统内部错误",
	CREATED:          "对象创建成功",
	ACCEPTED:         "请求已经被接受",
	NO_CONTENT:       "操作已经执行成功，但是没有返回数据",
	MOVED_PERM:       "资源已被移除",
	SEE_OTHER:        "重定向",
	NOT_MODIFIED:     "资源没有被修改",
	BAD_REQUEST:      "参数列表错误（缺少，格式不匹配）",
	UNAUTHORIZED:     "未授权",
	FORBIDDEN:        "访问受限，授权过期",
	NOT_FOUND:        "资源，服务未找到",
	BAD_METHOD:       "不允许的http方法",
	CONFLICT:         "资源冲突，或者资源被锁",
	UNSUPPORTED_TYPE: "不支持的数据，媒体类型",
	NOT_IMPLEMENTED:  "接口未实现",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
