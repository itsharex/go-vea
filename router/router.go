package router

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/constant"
	"go-vea/app/controller/common"
	"go-vea/app/controller/monitorctl"
	"go-vea/app/controller/systemctl"
	"go-vea/configs"
	"go-vea/middleware"
	"net/http"
)

func InitRouter() {
	r := gin.New()
	r.Use(gin.Recovery())
	// logrus 日志
	r.Use(middleware.Logger())
	// 静态文件
	r.StaticFS("uploads", http.Dir("uploads"))

	captcha := common.CaptchaHandler{}
	r.GET("/captchaImage", captcha.GetCaptcha)
	r.POST("/verify", captcha.VerifyCaptcha)

	sysLoginApi := systemctl.SysLoginApi{}
	r.POST("/login", sysLoginApi.Login)
	r.POST("/logout", sysLoginApi.Logout)
	r.GET("/getUserInfo", sysLoginApi.GetUserInfo)
	r.GET("/getRouters", sysLoginApi.GetRouters)

	/* 系统模块 */
	systemRoutes := r.Group("system")
	// jwt 认证
	systemRoutes.Use(middleware.JWT())

	// 文件处理
	fileRoutes := systemRoutes.Group("file")
	fileApi := common.FileApi{}
	{
		fileRoutes.POST("/upload", fileApi.UploadFile)
	}

	// 配置管理
	configRoutes := systemRoutes.Group("config")
	configApi := systemctl.SysConfigApi{}
	{
		configRoutes.POST("/list", middleware.HasPerm("system:config:list"), configApi.GetSysConfigList)
		configRoutes.GET("/:configId", middleware.HasPerm("system:config:query"), configApi.GetSysConfigById)
		configRoutes.GET("/configKey/:configKey", configApi.GetSysConfigByKey)
		configRoutes.POST("", middleware.OperationRecord("配置管理", constant.BUS_INSERT), middleware.HasPerm("system:config:add"), configApi.AddSysConfig)
		configRoutes.PUT("", middleware.OperationRecord("配置管理", constant.BUS_UPDATE), middleware.HasPerm("system:config:edit"), configApi.UpdateSysConfig)
		configRoutes.DELETE("", middleware.OperationRecord("配置管理", constant.BUS_DELETE), middleware.HasPerm("system:config:remove"), configApi.DeleteSysConfig)
		configRoutes.DELETE("/refreshCache", middleware.OperationRecord("配置管理", constant.BUS_CLEAN), configApi.RefreshCache)
	}

	// 字典管理
	dictRoutes := systemRoutes.Group("dict")
	dictDataApi := systemctl.SysDictDataApi{}
	dictTypeApi := systemctl.SysDictTypeApi{}
	{
		dictRoutes.POST("/data/list", middleware.HasPerm("system:dict:list"), dictDataApi.GetDictDataList)
		dictRoutes.GET("/data/type/:dictType", dictDataApi.GetDictDataListByDictType)
		dictRoutes.GET("/data/:dictCode", middleware.HasPerm("system:dict:query"), dictDataApi.GetDictData)
		dictRoutes.POST("/data", middleware.OperationRecord("字典数据", constant.BUS_INSERT), middleware.HasPerm("system:dict:add"), dictDataApi.AddDictData)
		dictRoutes.PUT("/data", middleware.OperationRecord("字典数据", constant.BUS_UPDATE), middleware.HasPerm("system:dict:update"), dictDataApi.UpdateDictData)
		dictRoutes.DELETE("/data", middleware.OperationRecord("字典数据", constant.BUS_DELETE), middleware.HasPerm("system:dict:remove"), dictDataApi.DeleteDictData)

		dictRoutes.POST("/type/list", middleware.HasPerm("system:dict:list"), dictTypeApi.GetDictTypeList)
		dictRoutes.GET("/type/:dictId", middleware.HasPerm("system:dict:query"), dictTypeApi.GetDictType)
		dictRoutes.POST("/type", middleware.OperationRecord("字典类型", constant.BUS_INSERT), middleware.HasPerm("system:dict:add"), dictTypeApi.AddDictType)
		dictRoutes.PUT("/type", middleware.OperationRecord("字典类型", constant.BUS_UPDATE), middleware.HasPerm("system:dict:edit"), dictTypeApi.UpdateDictType)
		dictRoutes.DELETE("/type", middleware.OperationRecord("字典类型", constant.BUS_DELETE), middleware.HasPerm("system:dict:remove"), dictTypeApi.DeleteDictType)
		dictRoutes.DELETE("/type/refreshCache", middleware.OperationRecord("字典类型", constant.BUS_CLEAN), middleware.HasPerm("system:dict:remove"), dictTypeApi.RefreshCache)
		dictRoutes.GET("/type/optionSelect", dictTypeApi.OptionSelect)
	}

	// 菜单管理
	menuRoutes := systemRoutes.Group("menu")
	menuApi := systemctl.SysMenuApi{}
	{
		menuRoutes.POST("", middleware.OperationRecord("菜单管理", constant.BUS_INSERT), middleware.HasPerm("system:menu:add"), menuApi.AddSysMenu)
		menuRoutes.PUT("", middleware.OperationRecord("菜单管理", constant.BUS_UPDATE), middleware.HasPerm("system:menu:edit"), menuApi.UpdateSysMenu)
		menuRoutes.DELETE("/:menuId", middleware.OperationRecord("菜单管理", constant.BUS_DELETE), middleware.HasPerm("system:menu:remove"), menuApi.DeleteSysMenu)
		menuRoutes.POST("/list", middleware.HasPerm("system:menu:list"), menuApi.GetMenuList)
		menuRoutes.POST("/listMenuTree", menuApi.GetMenuTreeList)
		menuRoutes.GET("/:menuId", middleware.HasPerm("system:menu:query"), menuApi.GetMenuInfo)
		menuRoutes.POST("/treeSelect", menuApi.TreeSelect)
		menuRoutes.POST("/roleMenuTreeSelect", menuApi.RoleMenuTreeSelect)
	}

	// 角色管理
	roleRoutes := systemRoutes.Group("role")
	roleApi := systemctl.SysRoleApi{}
	{
		roleRoutes.POST("/list", middleware.HasPerm("system:role:list"), roleApi.GetSysRoleList)
		roleRoutes.GET("/:roleId", middleware.HasPerm("system:role:query"), roleApi.GetSysRole)
		roleRoutes.POST("", middleware.OperationRecord("角色管理", constant.BUS_INSERT), middleware.HasPerm("system:role:add"), roleApi.AddSysRole)
		roleRoutes.PUT("", middleware.OperationRecord("角色管理", constant.BUS_UPDATE), middleware.HasPerm("system:role:edit"), roleApi.UpdateSysRole)
		roleRoutes.DELETE("", middleware.OperationRecord("角色管理", constant.BUS_DELETE), middleware.HasPerm("system:role:remove"), roleApi.DeleteSysRole)
		roleRoutes.PUT("dataScope", middleware.OperationRecord("角色管理", constant.BUS_UPDATE), middleware.HasPerm("system:role:edit"), roleApi.DataScope)
		roleRoutes.PUT("changeStatus", middleware.OperationRecord("角色管理", constant.BUS_UPDATE), middleware.HasPerm("system:role:edit"), roleApi.ChangeStatus)
		roleRoutes.GET("optionSelect", middleware.HasPerm("system:role:query"), roleApi.OptionSelect)
		roleRoutes.POST("/authUser/allocatedList", middleware.HasPerm("system:role:list"), roleApi.AllocatedList)
		roleRoutes.POST("/authUser/unallocatedList", middleware.HasPerm("system:role:list"), roleApi.UnallocatedList)
		roleRoutes.PUT("/authUser/cancel", middleware.OperationRecord("角色管理", constant.BUS_GRANT), middleware.HasPerm("system:role:edit"), roleApi.CancelAuthUser)
		roleRoutes.PUT("/authUser/cancelAll", middleware.OperationRecord("角色管理", constant.BUS_GRANT), middleware.HasPerm("system:role:edit"), roleApi.CancelAuthUserAll)
		roleRoutes.PUT("/authUser/selectAll", middleware.OperationRecord("角色管理", constant.BUS_GRANT), middleware.HasPerm("system:role:edit"), roleApi.SelectAuthUserAll)
		roleRoutes.GET("/deptTree/:roleId", middleware.HasPerm("system:role:query"), roleApi.DeptTree)
	}

	// 用户管理
	userRoutes := systemRoutes.Group("user")
	userApi := systemctl.SysUserApi{}
	profileApi := systemctl.SysProfileApi{}
	{
		userRoutes.POST("/list", middleware.HasPerm("system:user:list"), userApi.GetSysUserList)
		userRoutes.GET("/:userId", middleware.HasPerm("system:user:query"), userApi.GetSysUser)
		userRoutes.POST("", middleware.OperationRecord("用户管理", constant.BUS_INSERT), middleware.HasPerm("system:user:add"), userApi.AddSysUser)
		userRoutes.PUT("", middleware.OperationRecord("用户管理", constant.BUS_UPDATE), middleware.HasPerm("system:user:edit"), userApi.UpdateSysUser)
		userRoutes.DELETE("", middleware.OperationRecord("用户管理", constant.BUS_DELETE), middleware.HasPerm("system:user:delete"), userApi.DeleteSysUser)
		userRoutes.GET("deptTree", userApi.GetDeptTree)
		userRoutes.PUT("authRole", middleware.OperationRecord("用户管理", constant.BUS_GRANT), middleware.HasPerm("system:user:edit"), userApi.InsertAuthRole)
		userRoutes.GET("authRole/:userId", middleware.HasPerm("system:user:query"), userApi.AuthRole)
		userRoutes.PUT("resetPwd", middleware.OperationRecord("用户管理", constant.BUS_UPDATE), middleware.HasPerm("system:user:resetPwd"), userApi.ResetPwd)
		userRoutes.PUT("changeStatus", middleware.OperationRecord("用户管理", constant.BUS_UPDATE), middleware.HasPerm("system:user:edit"), userApi.ChangeStatus)

		userRoutes.GET("profile", profileApi.GetProfile)
		userRoutes.PUT("profile", middleware.OperationRecord("个人信息", constant.BUS_UPDATE), profileApi.UpdateProfile)
		userRoutes.PUT("/profile/updatePwd", middleware.OperationRecord("个人信息", constant.BUS_UPDATE), profileApi.UpdatePassword)
		userRoutes.POST("/profile/avatar", middleware.OperationRecord("用户头像", constant.BUS_UPDATE), profileApi.UploadAvatar)
	}

	// 通知管理
	noticeRoutes := systemRoutes.Group("notice")
	noticeApi := systemctl.SysNoticeApi{}
	{
		noticeRoutes.POST("/list", middleware.HasPerm("system:notice:list"), noticeApi.GetSysNoticeList)
		noticeRoutes.GET("/:noticeId", middleware.HasPerm("system:notice:query"), noticeApi.GetSysNotice)
		noticeRoutes.POST("", middleware.OperationRecord("通知管理", constant.BUS_INSERT), middleware.HasPerm("system:notice:add"), noticeApi.AddSysNotice)
		noticeRoutes.PUT("", middleware.OperationRecord("通知管理", constant.BUS_UPDATE), middleware.HasPerm("system:notice:edit"), noticeApi.UpdateSysNotice)
		noticeRoutes.DELETE("", middleware.OperationRecord("通知管理", constant.BUS_DELETE), middleware.HasPerm("system:notice:remove"), noticeApi.DeleteSysNotice)
	}

	// 部门管理
	deptRoutes := systemRoutes.Group("dept")
	deptApi := systemctl.SysDeptApi{}
	{
		deptRoutes.POST("/list", middleware.HasPerm("system:dept:list"), deptApi.GetSysDeptList)
		deptRoutes.POST("/tree", deptApi.GetSysDeptTreeList)
		deptRoutes.GET("/:deptId", middleware.HasPerm("system:dept:query"), deptApi.GetSysDept)
		deptRoutes.DELETE("/:deptId", middleware.OperationRecord("部门管理", constant.BUS_DELETE), middleware.HasPerm("system:dept:remove"), deptApi.DeleteSysDept)
		deptRoutes.POST("/list/exclude", middleware.HasPerm("system:dept:list"), deptApi.ExcludeChild)
		deptRoutes.POST("", middleware.OperationRecord("部门管理", constant.BUS_INSERT), middleware.HasPerm("system:dept:add"), deptApi.AddSysDept)
		deptRoutes.PUT("", middleware.OperationRecord("部门管理", constant.BUS_UPDATE), middleware.HasPerm("system:dept:edit"), deptApi.UpdateSysDept)
	}

	// 岗位管理
	postRoutes := systemRoutes.Group("post")
	postApi := systemctl.SysPostApi{}
	{
		postRoutes.POST("/list", middleware.HasPerm("system:post:list"), postApi.GetSysPostList)
		postRoutes.GET("/optionSelect", postApi.OptionSelect)
		postRoutes.GET("/:postId", middleware.HasPerm("system:post:query"), postApi.GetSysPost)
		postRoutes.POST("", middleware.OperationRecord("岗位管理", constant.BUS_INSERT), middleware.HasPerm("system:post:add"), postApi.AddSysPost)
		postRoutes.PUT("", middleware.OperationRecord("岗位管理", constant.BUS_UPDATE), middleware.HasPerm("system:post:edit"), postApi.UpdateSysPost)
		postRoutes.DELETE("", middleware.OperationRecord("岗位管理", constant.BUS_DELETE), middleware.HasPerm("system:post:remove"), postApi.DeleteSysPost)
	}

	/* 监控模块 */
	monitorRoutes := r.Group("monitor")
	// jwt 认证
	monitorRoutes.Use(middleware.JWT())

	// 系统信息
	serverInfoRouters := monitorRoutes.Group("server")
	serverApi := monitorctl.ServerInfoApi{}
	{
		serverInfoRouters.GET("info", middleware.HasPerm("monitor:server:list"), serverApi.GetServerInfo)
	}

	// 在线用户
	onlineUserRouters := monitorRoutes.Group("online")
	onlineUserApi := monitorctl.OnlineUserApi{}
	{
		onlineUserRouters.POST("list", middleware.HasPerm("monitor:online:list"), onlineUserApi.GetOnlineUser)
		onlineUserRouters.DELETE(":tokenId", middleware.OperationRecord("在线用户", constant.BUS_FORCE), middleware.HasPerm("monitor:online:forceLogout"), onlineUserApi.ForceLogout)
	}

	// 登录日志
	loginLogRoutes := monitorRoutes.Group("loginLog")
	loginLogApi := monitorctl.SysLoginLog{}
	{
		loginLogRoutes.POST("/list", middleware.HasPerm("monitor:loginLog:list"), loginLogApi.GetLoginLogList)
		loginLogRoutes.GET("/:id", middleware.HasPerm("monitor:loginLog:query"), loginLogApi.GetLoginLog)
		loginLogRoutes.DELETE("", middleware.OperationRecord("登录日志", constant.BUS_DELETE), middleware.HasPerm("monitor:loginLog:remove"), loginLogApi.DeleteLoginLog)
		loginLogRoutes.PUT("", middleware.OperationRecord("登录日志", constant.BUS_UPDATE), middleware.HasPerm("monitor:loginLog:edit"), loginLogApi.UpdateLoginLog)
		loginLogRoutes.DELETE("/clean", middleware.OperationRecord("登录日志", constant.BUS_CLEAN), middleware.HasPerm("monitor:loginLog:remove"), loginLogApi.CleanLoginLog)
		loginLogRoutes.GET("/unlock/:username", middleware.OperationRecord("登录日志", constant.BUS_OTHER), middleware.HasPerm("monitor:loginLog:unlock"), loginLogApi.Unlock)
	}

	// 操作日志
	operLogRoutes := monitorRoutes.Group("operLog")
	operLogApi := monitorctl.SysOperLog{}
	{
		operLogRoutes.POST("/list", middleware.HasPerm("monitor:operLog:list"), operLogApi.GetOperLogList)
		operLogRoutes.GET("/:operId", middleware.HasPerm("monitor:operLog:query"), operLogApi.GetOperLog)
		operLogRoutes.DELETE("", middleware.OperationRecord("操作日志", constant.BUS_DELETE), middleware.HasPerm("monitor:operLog:remove"), operLogApi.DeleteOperLog)
		operLogRoutes.DELETE("/clean", middleware.OperationRecord("操作日志", constant.BUS_CLEAN), middleware.HasPerm("monitor:operLog:remove"), operLogApi.CleanOperLog)
		operLogRoutes.PUT("", middleware.OperationRecord("操作日志", constant.BUS_UPDATE), middleware.HasPerm("monitor:operLog:edit"), operLogApi.UpdateOperLog)
	}

	_ = r.Run(configs.AppConfig.Server.Port)
}

func NewRouter() *gin.Engine {
	r := gin.New()
	return r
}
