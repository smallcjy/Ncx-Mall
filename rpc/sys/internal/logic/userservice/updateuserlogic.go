package userservicelogic

import (
	"context"
	"errors"
	"fmt"
	"github.com/feihua/zero-admin/rpc/sys/gen/model"
	"github.com/feihua/zero-admin/rpc/sys/gen/query"
	"github.com/feihua/zero-admin/rpc/sys/internal/svc"
	"github.com/feihua/zero-admin/rpc/sys/sysclient"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

// UpdateUserLogic 更新用户
/*
Author: LiuFeiHua
Date: 2023/12/18 14:37
*/
type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUser 更新用户(id为1是系统预留超级管理员用户,不能更新)
// 1.根据用户id查询用户是否已存在
// 2.查询用户名称是否存在
// 3.查询手机号是否存在
// 4.查询邮箱是否存在
// 5.用户不存在时,则直接添加用户
// 6.清空用户与岗位关联
// 7.添加用户与岗位关联
func (l *UpdateUserLogic) UpdateUser(in *sysclient.UpdateUserReq) (*sysclient.UpdateUserResp, error) {
	if in.Id == 1 {
		return nil, errors.New("不允许操作超级管理员用户")
	}

	name := in.UserName // 用户名

	user := query.SysUser
	q := user.WithContext(l.ctx)

	// 1.根据用户id查询用户是否已存在
	item, err := q.Where(query.SysUser.ID.Eq(in.Id)).First()

	// 1.判断用户是否存在
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errors.New("用户不存在")
	case err != nil:
		logc.Errorf(l.ctx, "查询用户异常, 请求参数：%+v, 异常信息: %s", in, err.Error())
		return nil, errors.New("查询用户异常")
	}

	// 2.查询用户名称是否存在
	count, err := q.Where(user.ID.Neq(in.Id), user.UserName.Eq(name)).Count()

	if err != nil {
		logc.Errorf(l.ctx, "根据用户名称：%s,查询用户失败,异常:%s", name, err.Error())
		return nil, errors.New(fmt.Sprintf("更新用户失败"))
	}

	if count > 0 {
		return nil, errors.New(fmt.Sprintf("用户：%s,已存在", name))
	}

	// 3.查询手机号是否存在
	mobile := in.Mobile
	count, err = q.Where(user.ID.Neq(in.Id), user.Mobile.Eq(mobile)).Count()

	if err != nil {
		logc.Errorf(l.ctx, "根据手机号：%s,查询用户失败,异常:%s", mobile, err.Error())
		return nil, errors.New(fmt.Sprintf("更新用户失败"))
	}

	if count > 0 {
		return nil, errors.New(fmt.Sprintf("手机号：%s,已存在", mobile))
	}

	// 4.查询邮箱是否存在
	email := in.Email
	count, err = q.Where(user.ID.Neq(in.Id), user.Email.Eq(email)).Count()

	if err != nil {
		logc.Errorf(l.ctx, "根据邮箱：%s,查询用户失败,异常:%s", email, err.Error())
		return nil, errors.New(fmt.Sprintf("更新用户失败"))
	}

	if count > 0 {
		return nil, errors.New(fmt.Sprintf("邮箱：%s,已存在", email))
	}

	// 5.用户存在时,则直接更新用户
	now := time.Now()
	sysUser := &model.SysUser{
		ID:           in.Id,             // 编号
		UserName:     in.UserName,       // 用户名
		NickName:     in.NickName,       // 昵称
		Avatar:       in.Avatar,         // 头像
		Password:     item.Password,     // 密码
		Salt:         in.Salt,           // 加密盐
		Email:        in.Email,          // 邮箱
		Mobile:       in.Mobile,         // 手机号
		UserStatus:   in.UserStatus,     // 帐号状态（1正常 0停用）
		DeptID:       in.DeptId,         // 部门id
		Remark:       in.Remark,         // 备注
		IsDeleted:    item.IsDeleted,    // 是否删除
		LoginTime:    item.LoginTime,    // 登录时间
		LoginIP:      item.LoginIP,      // 登录ip
		LoginOs:      item.LoginOs,      // 登录os
		LoginBrowser: item.LoginBrowser, // 登录浏览器
		CreateBy:     item.CreateBy,     // 创建者
		CreateTime:   item.CreateTime,   // 创建时间
		UpdateBy:     in.UpdateBy,       // 更新者
		UpdateTime:   &now,              // 更新时间
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		err = l.svcCtx.DB.Model(&model.SysUser{}).WithContext(l.ctx).Where(tx.SysUser.ID.Eq(sysUser.ID)).Save(sysUser).Error
		if err != nil {
			logc.Errorf(l.ctx, "更新用户异常,参数:%+v,异常:%s", user, err.Error())
			return err
		}

		var userPosts []*model.SysUserPost
		for _, postId := range in.PostIds {
			userPosts = append(userPosts, &model.SysUserPost{
				UserID: sysUser.ID,
				PostID: postId,
			})
		}

		postDo := tx.SysUserPost.WithContext(l.ctx)
		// 6.清空用户与岗位关联
		_, err = postDo.Where(tx.SysUserPost.UserID.Eq(sysUser.ID)).Delete()
		if err != nil {
			logc.Errorf(l.ctx, "删除用户与岗位关联异常,参数:%+v,异常:%s", user, err.Error())
			return err
		}

		// 7.添加用户与岗位关联
		err = postDo.CreateInBatches(userPosts, len(userPosts))
		if err != nil {
			logc.Errorf(l.ctx, "更新用户与岗位关联异常,参数:%+v,异常:%s", user, err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		logc.Errorf(l.ctx, "更新用户异常,参数:%+v,异常:%s", user, err.Error())
		return nil, errors.New("更新用户异常")
	}

	return &sysclient.UpdateUserResp{}, nil
}
