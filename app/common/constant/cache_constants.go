package constant

const (
	// LOGIN_TOKEN_KEY 登录用户 redis key
	LOGIN_TOKEN_KEY = "login_tokens:"

	// CAPTCHA_CODE_KEY 验证码 redis key
	CAPTCHA_CODE_KEY = "captcha_codes:"

	// SYS_CONFIG_KEY 参数管理 cache key
	SYS_CONFIG_KEY = "sys_config:"

	// SYS_DICT_KEY 字典管理 cache key
	SYS_DICT_KEY = "sys_dict:"

	// REPEAT_SUBMIT_KEY 防重提交 redis key
	REPEAT_SUBMIT_KEY = "repeat_submit:"

	// RATE_LIMIT_KEY 限流 redis key
	RATE_LIMIT_KEY = "rate_limit:"

	// PWD_ERR_CNT_KEY 登录账户密码错误次数 redis key
	PWD_ERR_CNT_KEY = "pwd_err_cnt:"
)
