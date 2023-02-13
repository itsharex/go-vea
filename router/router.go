package router

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/controller/common"
	"go-vea/app/controller/monitorctl"
	"go-vea/app/controller/systemctl"
	"go-vea/configs"
	"go-vea/middleware"
)

func InitRouter() {
	r := gin.New()
	r.Use(gin.Recovery())
	// logrus 日志
	r.Use(middleware.Logger())
	// 操作日志
	r.Use(middleware.OperationRecord())

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

	// 配置管理
	configRoutes := systemRoutes.Group("config")
	configApi := systemctl.SysConfigApi{}
	{
		configRoutes.POST("/list", configApi.GetSysConfigList, middleware.HasPerm("system:config:list"))
		configRoutes.GET("/:configId", configApi.GetSysConfigById)
		configRoutes.GET("/configKey/:configKey", configApi.GetSysConfigByKey)
		configRoutes.POST("", configApi.AddSysConfig)
		configRoutes.PUT("", configApi.UpdateSysConfig)
		configRoutes.DELETE("", configApi.DeleteSysConfig)
		configRoutes.DELETE("/refreshCache", configApi.RefreshCache)
	}

	// 字典管理
	dictRoutes := systemRoutes.Group("dict")
	dictDataApi := systemctl.SysDictDataApi{}
	dictTypeApi := systemctl.SysDictTypeApi{}
	{
		dictRoutes.POST("/data/list", dictDataApi.GetDictDataList)
		dictRoutes.GET("/data/type/:dictType", dictDataApi.GetDictDataListByDictType)
		dictRoutes.GET("/data/:dictCode", dictDataApi.GetDictData)
		dictRoutes.POST("/data", dictDataApi.AddDictData)
		dictRoutes.PUT("/data", dictDataApi.UpdateDictData)
		dictRoutes.DELETE("/data", dictDataApi.DeleteDictData)

		dictRoutes.POST("/type/list", dictTypeApi.GetDictTypeList)
		dictRoutes.GET("/type/:dictId", dictTypeApi.GetDictType)
		dictRoutes.POST("/type", dictTypeApi.AddDictType)
		dictRoutes.PUT("/type", dictTypeApi.UpdateDictType)
		dictRoutes.DELETE("/type", dictTypeApi.DeleteDictType)
		dictRoutes.DELETE("/type/refreshCache", dictTypeApi.RefreshCache)
		dictRoutes.GET("/type/optionSelect", dictTypeApi.OptionSelect)
	}

	// 菜单管理
	menuRoutes := systemRoutes.Group("menu")
	menuApi := systemctl.SysMenuApi{}
	{
		menuRoutes.POST("", menuApi.AddSysMenu)
		menuRoutes.PUT("", menuApi.UpdateSysMenu)
		menuRoutes.DELETE("/:menuId", menuApi.DeleteSysMenu)
		menuRoutes.POST("/list", menuApi.GetMenuList)
		menuRoutes.POST("/listMenuTree", menuApi.GetMenuTreeList)
		menuRoutes.GET("/:menuId", menuApi.GetMenuInfo)
		menuRoutes.POST("/treeSelect", menuApi.TreeSelect)
		menuRoutes.POST("/roleMenuTreeSelect", menuApi.RoleMenuTreeSelect)
	}

	// 角色管理
	roleRoutes := systemRoutes.Group("role")
	roleApi := systemctl.SysRoleApi{}
	{
		roleRoutes.POST("/list", roleApi.GetSysRoleList)
		roleRoutes.GET("/:roleId", roleApi.GetSysRole)
		roleRoutes.POST("", roleApi.AddSysRole)
		roleRoutes.PUT("", roleApi.UpdateSysRole)
		roleRoutes.DELETE("", roleApi.DeleteSysRole)
		roleRoutes.PUT("dataScope", roleApi.DataScope)
		roleRoutes.PUT("changeStatus", roleApi.ChangeStatus)
		roleRoutes.GET("optionSelect", roleApi.OptionSelect)
		roleRoutes.POST("/authUser/allocatedList", roleApi.AllocatedList)
		roleRoutes.POST("/authUser/unallocatedList", roleApi.UnallocatedList)
		roleRoutes.PUT("/authUser/cancel", roleApi.CancelAuthUser)
		roleRoutes.PUT("/authUser/cancelAll", roleApi.CancelAuthUserAll)
		roleRoutes.PUT("/authUser/selectAll", roleApi.SelectAuthUserAll)
		roleRoutes.GET("/deptTree/:roleId", roleApi.DeptTree)
	}

	// 用户管理
	userRoutes := systemRoutes.Group("user")
	userApi := systemctl.SysUserApi{}
	profileApi := systemctl.SysProfileApi{}
	{
		userRoutes.POST("/list", userApi.GetSysUserList)
		userRoutes.GET("/:userId", userApi.GetSysUser)
		userRoutes.POST("", userApi.AddSysUser)
		userRoutes.PUT("", userApi.UpdateSysUser)
		userRoutes.DELETE("", userApi.DeleteSysUser)
		userRoutes.GET("deptTree", userApi.GetDeptTree)
		userRoutes.PUT("authRole", userApi.InsertAuthRole)
		userRoutes.GET("authRole/:userId", userApi.AuthRole)
		userRoutes.PUT("resetPwd", userApi.ResetPwd)
		userRoutes.PUT("changeStatus", userApi.ChangeStatus)

		userRoutes.GET("profile", profileApi.GetProfile)
		userRoutes.PUT("profile", profileApi.UpdateProfile)
		userRoutes.PUT("/profile/updatePwd", profileApi.UpdatePassword)
		userRoutes.POST("/profile/avatar", profileApi.UploadAvatar)
	}

	// 通知管理
	noticeRoutes := systemRoutes.Group("notice")
	noticeApi := systemctl.SysNoticeApi{}
	{
		noticeRoutes.POST("/list", noticeApi.GetSysNoticeList)
		noticeRoutes.GET("/:noticeId", noticeApi.GetSysNotice)
		noticeRoutes.POST("", noticeApi.AddSysNotice)
		noticeRoutes.PUT("", noticeApi.UpdateSysNotice)
		noticeRoutes.DELETE("", noticeApi.DeleteSysNotice)
	}

	// 部门管理
	deptRoutes := systemRoutes.Group("dept")
	deptApi := systemctl.SysDeptApi{}
	{
		deptRoutes.POST("/list", deptApi.GetSysDeptList)
		deptRoutes.POST("/tree", deptApi.GetSysDeptTreeList)
		deptRoutes.GET("/:deptId", deptApi.GetSysDept)
		deptRoutes.DELETE("/:deptId", deptApi.DeleteSysDept)
		deptRoutes.POST("/list/exclude", deptApi.ExcludeChild)
		deptRoutes.POST("", deptApi.AddSysDept)
		deptRoutes.PUT("", deptApi.UpdateSysDept)
	}

	// 岗位管理
	postRoutes := systemRoutes.Group("post")
	postApi := systemctl.SysPostApi{}
	{
		postRoutes.POST("/list", postApi.GetSysPostList)
		postRoutes.GET("/optionSelect", postApi.OptionSelect)
		postRoutes.GET("/:postId", postApi.GetSysPost)
		postRoutes.POST("", postApi.AddSysPost)
		postRoutes.PUT("", postApi.UpdateSysPost)
		postRoutes.DELETE("", postApi.DeleteSysPost)
	}

	/* 监控模块 */
	monitorRoutes := r.Group("monitor")
	// jwt 认证
	monitorRoutes.Use(middleware.JWT())

	// 系统信息
	serverInfoRouters := monitorRoutes.Group("server")
	serverApi := monitorctl.ServerInfoApi{}
	{
		serverInfoRouters.GET("info", serverApi.GetServerInfo)
	}

	// 在线用户
	onlineUserRouters := monitorRoutes.Group("online")
	onlineUserApi := monitorctl.OnlineUserApi{}
	{
		onlineUserRouters.POST("list", onlineUserApi.GetOnlineUser)
		onlineUserRouters.DELETE(":tokenId", onlineUserApi.ForceLogout)
	}

	// 登录日志
	loginLogRoutes := monitorRoutes.Group("loginLog")
	loginLogApi := monitorctl.SysLoginLog{}
	{
		loginLogRoutes.POST("/list", loginLogApi.GetLoginLogList)
		loginLogRoutes.GET("/:id", loginLogApi.GetLoginLog)
		loginLogRoutes.DELETE("", loginLogApi.DeleteLoginLog)
		loginLogRoutes.PUT("", loginLogApi.UpdateLoginLog)
		loginLogRoutes.DELETE("/clean", loginLogApi.CleanLoginLog)
		loginLogRoutes.GET("/unlock/:username", loginLogApi.Unlock)
	}

	// 操作日志
	operLogRoutes := monitorRoutes.Group("operLog")
	operLogApi := monitorctl.SysOperLog{}
	{
		operLogRoutes.POST("/list", operLogApi.GetOperLogList)
		operLogRoutes.GET("/:operId", operLogApi.GetOperLog)
		operLogRoutes.DELETE("", operLogApi.DeleteOperLog)
		operLogRoutes.DELETE("/clean", operLogApi.CleanOperLog)
		operLogRoutes.PUT("", operLogApi.UpdateOperLog)
	}

	_ = r.Run(configs.AppConfig.Server.Port)
}

func NewRouter() *gin.Engine {
	r := gin.New()
	return r
}
