package syssrv

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
	"go-vea/util"
	"gorm.io/gorm"
	"strings"
)

type SysUserService struct{}

var SysUserSrv = new(SysUserService)

func (*SysUserService) GetSysUserList(ctx context.Context, sysUser *request.SysUser) (*page.Pagination, error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.SelectList(sysUser)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysUserService) GetSysUserById(ctx context.Context, roleId int64) (*system.SysUser, error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.SelectById(roleId)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysUserService) AddSysUser(ctx context.Context, addSysUser *request.AddSysUser) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	// Transactional begin
	tx := sysUserDao.Begin()
	// 新增用户信息
	addSysUser.SysUser.Password, _ = util.PasswordHash(addSysUser.SysUser.Password)
	err := sysUserDao.Insert(addSysUser.SysUser)
	// 新增用户岗位关联
	err = addUserPost(ctx, addSysUser)
	// 新增用户与角色管理
	err = addUserRole(ctx, addSysUser)
	if err != nil {
		tx.Rollback()
		global.Logger.Error(err)
		return err
	}
	// Transactional commit
	tx.Commit()
	return nil
}

func (*SysUserService) UpdateUserById(ctx context.Context, addSysUser *request.AddSysUser) (err error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	sysUserRoleDao := sysdao.NewSysUserRoleDao(ctx)
	sysUserPostDao := sysdao.NewSysUserPostDao(ctx)
	// todo Transactional
	// 删除用户与角色关联
	err = sysUserRoleDao.DeleteUserRoleByUserId(addSysUser.SysUser.UserID)
	// 新增用户与角色管理
	err = addUserRole(ctx, addSysUser)
	// 删除用户与岗位关联
	err = sysUserPostDao.DeleteUserPostByUserId(addSysUser.SysUser.UserID)
	// 新增用户与岗位管理
	err = addUserPost(ctx, addSysUser)
	// 更新用户信息
	err = sysUserDao.UpdateById(addSysUser.SysUser)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) DeleteSysUserById(ctx context.Context, userId int64) (err error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	sysUserRoleDao := sysdao.NewSysUserRoleDao(ctx)
	sysUserPostDao := sysdao.NewSysUserPostDao(ctx)
	// todo Transactional
	// 删除用户与角色关联
	err = sysUserRoleDao.DeleteUserRoleByUserId(userId)
	// 删除用户与岗位表
	err = sysUserPostDao.DeleteUserPostByUserId(userId)
	err = sysUserDao.DeleteById(userId)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) DeleteSysUserByIds(ctx context.Context, ids []int64) (err error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	sysUserRoleDao := sysdao.NewSysUserRoleDao(ctx)
	sysUserPostDao := sysdao.NewSysUserPostDao(ctx)
	// todo Transactional
	// 删除用户与角色关联
	err = sysUserRoleDao.DeleteUserRole(ids)
	// 删除用户与岗位表
	err = sysUserPostDao.DeleteUserPost(ids)
	err = sysUserDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) SelectUserRoleGroup(ctx context.Context, username string) (string, error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	roles, err := sysRoleDao.SelectRolesByUserName(username)
	if err != nil {
		return "", err
	}
	if roles != nil {
		var roleNames []string
		for _, role := range roles {
			roleNames = append(roleNames, role.RoleName)
		}
		res := strings.Join(roleNames, ",")
		return res, nil
	}
	return "", nil
}

func (*SysUserService) SelectUserPostGroup(ctx context.Context, username string) (string, error) {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	posts, err := sysPostDao.SelectPostsByUserName(username)
	if err != nil {
		return "", err
	}
	if posts != nil {
		var postNames []string
		for _, post := range posts {
			postNames = append(postNames, post.PostName)
		}
		res := strings.Join(postNames, ",")
		return res, nil
	}
	return "", nil
}

func (*SysUserService) ResetPwd(ctx context.Context, sysUser *system.SysUser) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	err := sysUserDao.UpdateById(sysUser)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) InsertUserAuth(ctx context.Context, params *request.AddUserRole) (err error) {
	sysUserRoleDao := sysdao.NewSysUserRoleDao(ctx)
	// todo Transactional
	err = sysUserRoleDao.DeleteUserRoleByUserId(params.UserId)
	err = insertUserRole(ctx, params)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) UpdateUserStatus(ctx *gin.Context, sysUser *system.SysUser) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	err := sysUserDao.UpdateById(sysUser)
	return err
}

func (*SysUserService) CheckUserNameUnique(ctx context.Context, sysUser *system.SysUser) bool {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.CheckUserNameUnique(sysUser.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.UserID != 0 && data.UserID != sysUser.UserID {
		global.Logger.Error("")
		return false
	}
	return true
}

func (*SysUserService) CheckPhoneUnique(ctx context.Context, sysUser *system.SysUser) bool {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.CheckPhoneUnique(sysUser.PhoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.UserID != 0 && data.UserID != sysUser.UserID {
		return false
	}
	return true
}

func (*SysUserService) CheckEmailUnique(ctx context.Context, sysUser *system.SysUser) bool {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.CheckEmailUnique(sysUser.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.UserID != 0 && data.UserID != sysUser.UserID {
		return false
	}
	return true
}

func (*SysUserService) CheckUserAllowed(ctx *gin.Context, sysUser *system.SysUser) bool {
	//loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
	//if sysUser != nil && loginUser.SysUserResp.SysUser.IsAdmin(sysUser.UserID) {
	//	global.Logger.Error("不允许操作超级管理员用户")
	//	return false
	//}
	return true
}

func (s *SysUserService) CheckUserDataScope(ctx *gin.Context, userId int64) bool {
	//loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
	//if !loginUser.SysUserResp.SysUser.IsAdmin(userId) {
	//	params := &request.SysUser{
	//		UserID: userId,
	//	}
	//	data, _ := s.GetSysUserList(ctx, params)
	//	if data.Rows == nil {
	//		global.Logger.Error("没有权限访问用户数据")
	//		return false
	//	}
	//}
	return true
}

func addUserRole(ctx context.Context, params *request.AddSysUser) error {
	userRole := &request.AddUserRole{
		UserId:  params.SysUser.UserID,
		RoleIds: params.RoleIds,
	}
	err := insertUserRole(ctx, userRole)
	return err
}

func addUserPost(ctx context.Context, params *request.AddSysUser) error {
	userPost := &request.AddUserPost{
		UserId:  params.SysUser.UserID,
		PostIds: params.PostIds,
	}
	err := insertUserPost(ctx, userPost)
	return err
}

func insertUserRole(ctx context.Context, params *request.AddUserRole) (err error) {
	sysUserRoleDao := sysdao.NewSysUserRoleDao(ctx)
	if params.RoleIds != nil {
		var sysUserRoleList []*system.SysUserRole
		for _, roleId := range params.RoleIds {
			sysUserRole := &system.SysUserRole{
				UserID: params.UserId,
				RoleID: roleId,
			}
			sysUserRoleList = append(sysUserRoleList, sysUserRole)
		}
		err = sysUserRoleDao.BatchUserRole(sysUserRoleList)
	}
	return err
}

func insertUserPost(ctx context.Context, params *request.AddUserPost) (err error) {
	sysUserPostDao := sysdao.NewSysUserPostDao(ctx)
	if params.PostIds != nil {
		var sysUserPostList []*system.SysUserPost
		for _, postId := range params.PostIds {
			sysUserPost := &system.SysUserPost{
				UserID: params.UserId,
				PostID: postId,
			}
			sysUserPostList = append(sysUserPostList, sysUserPost)
		}
		err = sysUserPostDao.BatchUserPost(sysUserPostList)
	}
	return err
}
