package sysdao

import (
	"context"
	"go-vea/app/common/e"
	"go-vea/app/common/page"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysMenuDao struct {
	*gorm.DB
}

func NewSysMenuDao(ctx context.Context) *SysMenuDao {
	return &SysMenuDao{configs.GetDB(ctx)}
}

func NewSysMenuDaoByDB(db *gorm.DB) *SysMenuDao {
	return &SysMenuDao{db}
}

func (dao *SysMenuDao) SelectList(sysMenu *request.SysMenu) (p *page.Pagination, err error) {
	var menuList []*system.SysMenu
	p = new(page.Pagination)

	if sysMenu.ParentID != 0 {
		dao.DB = dao.DB.Where("parent_id = ?", sysMenu.ParentID)
	}
	if sysMenu.MenuName != "" {
		dao.DB = dao.DB.Where("menu_name = ?", sysMenu.MenuName)
	}
	if sysMenu.Visible != "" {
		dao.DB = dao.DB.Where("visible = ?", sysMenu.Visible)
	}
	if sysMenu.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysMenu.Status)
	}

	if sysMenu.OpenPage {
		p.PageNum = sysMenu.PageNum
		p.PageSize = sysMenu.PageSize
		err = dao.DB.Scopes(page.SelectPage(menuList, p, dao.DB)).Find(&menuList).Error
	} else {
		err = dao.DB.Find(&menuList).Error
	}
	p.Rows = menuList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysMenuDao) SelectById(menuId int64) (sysMenu *system.SysMenu, err error) {
	err = dao.DB.Where("menu_id = ?", menuId).Find(&sysMenu).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *SysMenuDao) Insert(sysMenu *system.SysMenu) error {
	return dao.DB.Create(sysMenu).Error
}

func (dao *SysMenuDao) UpdateById(sysMenu *system.SysMenu) error {
	return dao.DB.Save(sysMenu).Error
}

func (dao *SysMenuDao) DeleteById(menuId int64) error {
	return dao.DB.Where("menu_id = ?", menuId).Delete(&system.SysMenu{}).Error
}

func (dao *SysMenuDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("menu_id in (?)", ids).Delete(&system.SysMenu{}).Error
}

// CheckMenuNameUnique 校验菜单名称是否唯一
func (dao *SysMenuDao) CheckMenuNameUnique(sysMenu *system.SysMenu) (count int64, err error) {
	err = dao.DB.Model(&system.SysMenu{}).Where("menu_name = ?", sysMenu.MenuName).Where("parent_id = ?", sysMenu.ParentID).Count(&count).Error
	return
}

// HasChildByMenuId 是否存在菜单子节点
func (dao *SysMenuDao) HasChildByMenuId(menuId int64) (count int64, err error) {
	err = dao.DB.Model(&system.SysMenu{}).Where("parent_id = ?", menuId).Count(&count).Error
	return
}

// CheckMenuExistRole 查询菜单使用数量
func (dao *SysMenuDao) CheckMenuExistRole(menuId int64) (count int64, err error) {
	err = dao.DB.Model(&system.SysMenu{}).Where("menu_id = ?", menuId).Count(&count).Error
	return
}

// SelectMenuPerms 根据用户所有权限
func (dao *SysMenuDao) SelectMenuPerms() (list []string, err error) {
	// 		select distinct m.perms
	//		from sys_menu m
	//			 left join sys_role_menu rm on m.menu_id = rm.menu_id
	//			 left join sys_user_role ur on rm.role_id = ur.role_id
	err = dao.DB.Table("sys_menu m").Distinct("m.perms").Select("m.perms").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Joins("left join sys_user_role ur on rm.role_id = ur.role_id").
		Find(&list).Error
	return list, err
}

// SelectMenuPermsByUserId 根据用户ID查询权限
func (dao *SysMenuDao) SelectMenuPermsByUserId(userId int64) (list []string, err error) {
	// 		select distinct m.perms
	//		from sys_menu m
	//			 left join sys_role_menu rm on m.menu_id = rm.menu_id
	//			 left join sys_user_role ur on rm.role_id = ur.role_id
	//			 left join sys_role r on r.role_id = ur.role_id
	//		where m.status = '0' and r.status = '0' and ur.user_id = #{userId}
	err = dao.DB.Table("sys_menu m").Distinct("m.perms").Select("m.perms").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Joins("left join sys_user_role ur on rm.role_id = ur.role_id").
		Joins("left join sys_role r on r.role_id = ur.role_id").
		Where("m.status = '0' and r.status = '0' and ur.user_id = ?", userId).
		Find(&list).Error
	return list, err
}

// SelectMenuPermsByRoleId 根据角色ID查询权限
func (dao *SysMenuDao) SelectMenuPermsByRoleId(roleId int64) (list []string, err error) {
	// 		select distinct m.perms
	//		from sys_menu m
	//			 left join sys_role_menu rm on m.menu_id = rm.menu_id
	//		where m.status = '0' and rm.role_id = #{roleId}
	var perms []string
	err = dao.DB.Table("sys_menu m").Select("m.menu_id").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Where("rm.role_id = ?", roleId).
		Order("m.parent_id, m.order_num").
		Find(&perms).
		Error
	return perms, err
}

// SelectMenuTreeAll 查询系统菜单列表 M、C
func (dao *SysMenuDao) SelectMenuTreeAll() (menus []*system.SysMenu, err error) {
	err = dao.DB.Model(&system.SysMenu{}).Where("menu_type in ('M','C')").Where("status = ?", 0).Order("parent_id,order_num").Find(&menus).Error
	return
}

// SelectMenuTreeByUserId 根据用户Id查询系统菜单列表 M、C 涉及sys_role_menu、sys_user_role、sys_role、sys_user
func (dao *SysMenuDao) SelectMenuTreeByUserId(userId int64) (menus []*system.SysMenu, err error) {
	// 		select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time
	//		from sys_menu m
	//			 left join sys_role_menu rm on m.menu_id = rm.menu_id
	//			 left join sys_user_role ur on rm.role_id = ur.role_id
	//			 left join sys_role ro on ur.role_id = ro.role_id
	//			 left join sys_user u on ur.user_id = u.user_id
	//		where u.user_id = #{userId} and m.menu_type in ('M', 'C') and m.status = 0  AND ro.status = 0
	//		order by m.parent_id, m.order_num
	err = dao.DB.Table("sys_menu m").
		Select("m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Joins("left join sys_user_role ur on rm.role_id = ur.role_id").
		Joins("left join sys_role ro on ur.role_id = ro.role_id").
		Joins("left join sys_user u on ur.user_id = u.user_id").
		Where("u.user_id = ?", userId).
		Where("m.menu_type in ('M','C')").
		Where("m.status = 0 and ro.status = 0").
		Order("m.parent_id,m.order_num").
		Find(&menus).Error
	return
}

// SelectMenuListByRoleId 根据角色ID查询菜单树信息
func (dao *SysMenuDao) SelectMenuListByRoleId(param *request.MenuListByRoleId) (list []int64, err error) {
	// 		select m.menu_id
	//		from sys_menu m
	//            left join sys_role_menu rm on m.menu_id = rm.menu_id
	//        where rm.role_id = #{roleId}
	//            <if test="menuCheckStrictly">
	//              and m.menu_id not in (select m.parent_id from sys_menu m inner join sys_role_menu rm on m.menu_id = rm.menu_id and rm.role_id = #{roleId})
	//            </if>
	//		order by m.parent_id, m.order_num
	dao.DB = dao.DB.Table("sys_menu m").Distinct("m.perms").Select("m.perms").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id")

	if param.IsMenuCheckStrictly {
		dao.DB = dao.DB.Where("m.menu_id not in (select m.parent_id from sys_menu m inner join sys_role_menu rm on m.menu_id = rm.menu_id and rm.role_id = ?)", param.RoleId)
	}

	err = dao.DB.Where("m.status = '0' and rm.role_id = ?", param.RoleId).Find(&list).Error
	return list, err
}

// SelectMenuListByUserId 根据用户查询系统菜单列表
func (dao *SysMenuDao) SelectMenuListByUserId(sysMenu *request.SysMenu) (p *page.Pagination, err error) {
	var menuList []*system.SysMenu
	p = new(page.Pagination)

	err = dao.DB.Table("sys_menu m").
		Select("m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Joins("left join sys_user_role ur on rm.role_id = ur.role_id").
		Joins("left join sys_role ro on ur.role_id = ro.role_id").
		Where("u.user_id = ?", sysMenu.UserId).
		Find(&menuList).Error

	if sysMenu.ParentID != 0 {
		dao.DB = dao.DB.Where("parent_id = ?", sysMenu.ParentID)
	}
	if sysMenu.MenuName != "" {
		dao.DB = dao.DB.Where("menu_name = ?", sysMenu.MenuName)
	}
	if sysMenu.Visible != "" {
		dao.DB = dao.DB.Where("visible = ?", sysMenu.Visible)
	}
	if sysMenu.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysMenu.Status)
	}

	p.Rows = menuList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}
