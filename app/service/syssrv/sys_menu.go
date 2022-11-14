package syssrv

import (
	"context"
	"errors"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/global"
	"go-vea/util"
	"strings"
)

type SysMenuService struct{}

var SysMenuSrv = new(SysMenuService)

func (*SysMenuService) SelectMenuList(ctx context.Context, sysMenu *request.SysMenu, userId int64) (*page.Pagination, error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	sysUser := system.SysUser{}
	sysMenu.UserId = userId
	if sysUser.IsAdmin(userId) {
		data, err := sysMenuDao.SelectList(sysMenu)
		if err != nil {
			global.Logger.Error(err)
			return nil, err
		}
		return data, nil
	} else {
		data, err := sysMenuDao.SelectMenuListByUserId(sysMenu)
		if err != nil {
			global.Logger.Error(err)
			return nil, err
		}
		return data, nil
	}
}

func (*SysMenuService) SelectSysMenuById(ctx context.Context, menuId int64) (*system.SysMenu, error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	data, err := sysMenuDao.SelectById(menuId)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

func (s *SysMenuService) AddSysMenu(ctx context.Context, sysMenu *system.SysMenu) error {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	r := checkMenuNameUnique(ctx, sysMenu)
	if r {
		global.Logger.Error("新增失败！菜单名称已存在")
		return errors.New("新增失败！菜单名称已存在")
	} else if sysMenu.IsFrame == 0 && !isHttp(sysMenu.Path) {
		global.Logger.Error("地址必须以http(s)://开头")
		return errors.New("地址必须以http(s)://开头")
	}
	err := sysMenuDao.Insert(sysMenu)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysMenuService) UpdateSysMenuById(ctx context.Context, sysMenu *system.SysMenu) error {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	hasMenuName := checkMenuNameUnique(ctx, sysMenu)
	if hasMenuName {
		global.Logger.Error("修改失败！菜单名称已存在")
		return errors.New("修改失败！菜单名称已存在")
	} else if sysMenu.IsFrame == 0 && !isHttp(sysMenu.Path) {
		global.Logger.Error("地址必须以http(s)://开头")
		return errors.New("地址必须以http(s)://开头")
	} else if sysMenu.MenuID == sysMenu.ParentID {
		global.Logger.Error("上级菜单不能选择自己")
		return errors.New("上级菜单不能选择自己")
	}
	err := sysMenuDao.UpdateById(sysMenu)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysMenuService) DeleteSysMenuByIds(ctx context.Context, menuId int64) error {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	hasChild := hasChildByMenuId(ctx, menuId)
	existRole := checkMenuExistRole(ctx, menuId)
	if hasChild {
		global.Logger.Error("存在子菜单,不允许删除")
		return errors.New("存在子菜单,不允许删除")
	}
	if existRole {
		global.Logger.Error("菜单已分配,不允许删除")
		return errors.New("菜单已分配,不允许删除")
	}
	err := sysMenuDao.DeleteById(menuId)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysMenuService) selectMenuPermsByUserId(ctx context.Context, userId int64) ([]string, error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	listPerms, err := sysMenuDao.SelectMenuPermsByUserId(userId)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	var permsSet []string
	for _, perm := range listPerms {
		if perm != "" {
			permsSet = append(permsSet, perm)
		}
	}
	return permsSet, nil
}

func (*SysMenuService) selectMenuPermsByRoleId(ctx context.Context, roleId int64) ([]string, error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	listPerms, err := sysMenuDao.SelectMenuPermsByRoleId(roleId)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	var permsSet []string
	for _, perm := range listPerms {
		if perm != "" {
			permsSet = append(permsSet, perm)
		}
	}
	return permsSet, nil
}

// SelectMenuTreeByUserId 根据用户ID查询菜单
func (*SysMenuService) SelectMenuTreeByUserId(ctx context.Context, sysUser *system.SysUser) (menus []*system.SysMenu, err error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	if sysUser.IsAdmin(sysUser.UserID) {
		menus, err = sysMenuDao.SelectMenuTreeAll()
	} else {
		menus, err = sysMenuDao.SelectMenuTreeByUserId(sysUser.UserID)
	}
	return buildMenuTree(menus), nil
}

func (*SysMenuService) SelectMenuListByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	sysRole, err := sysRoleDao.SelectById(roleId)
	menuListByRoleId := &request.MenuListByRoleId{
		RoleId:              sysRole.RoleID,
		IsMenuCheckStrictly: sysRole.MenuCheckStrictly,
	}
	list, err := sysMenuDao.SelectMenuListByRoleId(menuListByRoleId)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return list, nil
}

func (*SysMenuService) GetBuildMenus(menus []*system.SysMenu) []response.RouterVo {
	return buildMenus(menus)
}

func (s *SysMenuService) GetTreeSelect(ctx context.Context, menu *request.SysMenu, userId int64) ([]*system.SysMenu, error) {
	data, err := s.SelectMenuList(ctx, menu, userId)
	if err != nil {
		return nil, err
	}
	// 类型转换
	menus, ok := data.Rows.([]*system.SysMenu)
	if !ok {
		global.Logger.Error("类型转换错误")
		return nil, errors.New("类型转换错误")
	}
	treeList := buildMenuTree(menus)
	return treeList, nil
}

// 是否存在菜单子节点
func hasChildByMenuId(ctx context.Context, menuId int64) bool {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	c, e := sysMenuDao.HasChildByMenuId(menuId)
	if e != nil {
		global.Logger.Error(e)
		return false
	}
	return c > 0
}

// 校验菜单名称是否唯一
func checkMenuNameUnique(ctx context.Context, menu *system.SysMenu) bool {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	c, e := sysMenuDao.CheckMenuNameUnique(menu)
	if e != nil {
		global.Logger.Error(e)
		return false
	}
	return c > 0
}

// 查询菜单使用数量
func checkMenuExistRole(ctx context.Context, menuId int64) bool {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	c, e := sysMenuDao.CheckMenuExistRole(menuId)
	if e != nil {
		global.Logger.Error(e)
		return false
	}
	return c > 0
}

// 构建菜单树
func buildMenuTree(menus []*system.SysMenu) []*system.SysMenu {
	menuMap := make(map[int64]*system.SysMenu)
	for i, menu := range menus {
		menu.ArrIdx = i
		menuMap[menu.MenuID] = menu
	}

	var resList []*system.SysMenu
	var childList []*system.SysMenu

	for _, menu := range menus {

		parent, ok := menuMap[menu.ParentID]

		if ok {
			if len(parent.Children) == 0 {
				childList = menus[parent.ArrIdx].Children
				if childList == nil {
					childList = []*system.SysMenu{}
				}
				childList = append(childList, menu)
				menus[parent.ArrIdx].Children = childList
			} else {
				menus[parent.ArrIdx].Children = append(menus[parent.ArrIdx].Children, menu)
			}
		}

		if menu.ParentID == 0 {
			resList = append(resList, menu)
		}
	}
	return resList
}

// BuildMenus 构建前端路由所需要的菜单
func buildMenus(menus []*system.SysMenu) []response.RouterVo {
	var routers []response.RouterVo

	for _, menu := range menus {
		router := response.RouterVo{
			Hidden:    menu.Visible == "1",
			Name:      getRouteName(menu),
			Path:      getRouterPath(menu),
			Component: getComponent(menu),
			Query:     menu.Query,
			Meta: &response.MetaVo{
				Title:   menu.MenuName,
				Icon:    menu.Icon,
				NoCache: menu.IsCache == 1,
				Link:    setMetaLink(menu.Path),
			},
		}

		childMenus := menu.Children
		if len(childMenus) > 0 && childMenus != nil && menu.MenuType == "M" {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
			router.Children = buildMenus(childMenus)
		} else if isMenuFrame(menu) {
			var m *response.MetaVo = nil
			router.Meta = m

			var childrenList []response.RouterVo
			var child response.RouterVo
			child.Path = menu.Path
			child.Component = menu.Component
			routerName := util.FirstUpper(menu.MenuName)
			child.Name = routerName
			child.Query = menu.Query
			child.Meta = &response.MetaVo{
				Title:   menu.MenuName,
				Icon:    menu.Icon,
				NoCache: menu.IsCache == 1,
				Link:    setMetaLink(menu.Path),
			}
			childrenList = append(childrenList, child)
			router.Children = childrenList
		} else if menu.ParentID == 0 && isInnerLink(menu) {
			router.Meta = &response.MetaVo{
				Title: menu.MenuName,
				Icon:  menu.Icon,
			}
			router.Path = menu.Path

			var childrenList []response.RouterVo
			var child response.RouterVo

			routerPath := innerLinkReplaceEach(menu.Path)
			child.Path = routerPath
			child.Component = "InnerLink"
			routerName := util.FirstUpper(menu.MenuName)
			child.Name = routerName
			child.Query = menu.Query
			child.Meta = &response.MetaVo{
				Title: menu.MenuName,
				Icon:  menu.Icon,
				Link:  menu.Path,
			}
			childrenList = append(childrenList, child)
			router.Children = childrenList
		}
		routers = append(routers, router)
	}

	return routers
}

// 获取路由名称
func getRouteName(sysMenu *system.SysMenu) string {
	routerName := util.FirstUpper(sysMenu.Path)
	// 非外链并且是一级目录（类型为目录）
	if isMenuFrame(sysMenu) {
		routerName = ""
	}
	return routerName
}

// 获取路由地址
func getRouterPath(sysMenu *system.SysMenu) string {
	routerPath := sysMenu.Path
	// 内链打开外网方式
	if sysMenu.ParentID != 0 && isInnerLink(sysMenu) {
		routerPath = innerLinkReplaceEach(routerPath)
	}
	// 非外链并且是一级目录（类型为目录）
	if sysMenu.ParentID == 0 && sysMenu.MenuType == "M" && sysMenu.IsFrame == 1 {
		routerPath = "/" + sysMenu.Path
	} else if isMenuFrame(sysMenu) {
		// 非外链并且是一级目录（类型为菜单）
		routerPath = "/"
	}
	return routerPath
}

// 获取组件信息
func getComponent(sysMenu *system.SysMenu) string {
	component := "Layout"
	if sysMenu.Component != "" && !isMenuFrame(sysMenu) {
		component = sysMenu.Component
	} else if sysMenu.Component == "" && sysMenu.ParentID != 0 && isInnerLink(sysMenu) {
		component = "InnerLink"
	} else if sysMenu.Component == "" && isParentView(sysMenu) {
		component = "ParentView"
	}
	return component
}

// 是否为菜单内部跳转
func isMenuFrame(sysMenu *system.SysMenu) bool {
	return sysMenu.ParentID == 0 && sysMenu.MenuType == "C" && sysMenu.IsFrame == 1
}

// 是否为内链组件
func isInnerLink(sysMenu *system.SysMenu) bool {
	return sysMenu.IsFrame == 1 && isHttp(sysMenu.Path)
}

// 是否为parent_view组件
func isParentView(sysMenu *system.SysMenu) bool {
	return sysMenu.ParentID != 0 && sysMenu.MenuType == "M"
}

// 内链域名特殊字符替换
func innerLinkReplaceEach(str string) string {
	str1 := strings.Replace(str, "http://", "", 1)
	str2 := strings.Replace(str1, "https://", "", 1)
	str3 := strings.Replace(str2, "www.", "", 1)
	str4 := strings.Replace(str3, ".", "/", 1)
	return str4
}

func isHttp(path string) bool {
	return strings.HasPrefix(path, "http://")
}

func setMetaLink(path string) string {
	if isHttp(path) {
		return path
	}
	return ""
}
