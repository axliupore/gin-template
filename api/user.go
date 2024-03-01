package api

import (
	"errors"
	"github.com/axliupore/gin-template/global"
	"github.com/axliupore/gin-template/model"
	"github.com/axliupore/gin-template/model/request"
	"github.com/axliupore/gin-template/model/response"
	"github.com/axliupore/gin-template/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// UserRegister
// @Tags     user
// @Summary  用户注册
// @accept   application/json
// @Produce  application/json
// @Param    data  body      request.UserRegister      				  true  "账号,密码,确认密码"
// @Success  200   {object}  response.Response{data=model.User,msg=string}  "返回用户信息"
// @Router   /api/user/register [post]
func UserRegister(c *gin.Context) {
	var r request.UserRegister
	err := c.ShouldBindJSON(&r)
	if err != nil || utils.IsAnyBlank(r.Account, r.Password, r.CheckPassword) {
		response.Error(c, utils.Params)
		return
	}
	if len(r.Account) < 4 {
		response.ErrorMessage(c, utils.Params, "账号不能小于4位")
		return
	}
	if len(r.Password) < 8 || len(r.CheckPassword) < 8 || r.Password != r.CheckPassword {
		response.ErrorMessage(c, utils.Params, "密码不能小于8位或两次密码不一致")
		return
	}
	user := model.NewUser(r.Account, r.Password)
	err = userService.UserRegister(user)
	if err != nil {
		global.Log.Error("注册失败!", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorRegister, err.Error())
		return
	}
	response.SuccessData(c, *user)
}

// UserLogin
// @Tags     user
// @Summary  用户登录
// @accept   application/json
// @Produce  application/json
// @Param    data  body      request.UserLogin      				                  true  "账号,密码"
// @Success  200   {object}  response.Response{data=response.UserLoginResponse,msg=string}  "返回用户信息,token,过期时间"
// @Router   /api/user/login [post]
func UserLogin(c *gin.Context) {
	var r request.UserLogin
	err := c.ShouldBindJSON(&r)
	if err != nil || utils.IsAnyBlank(r.Account, r.Password) {
		response.Error(c, utils.Params)
		return
	}
	user := model.NewUser(r.Account, r.Password)
	// 查询用户信息
	user, err = userService.UserLogin(user)
	if err != nil {
		global.Log.Error("登录失败！用户不存在或密码错误", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorLogin, "用户名不存在或者密码错误")
		return
	}
	if user.IsDelete != 0 {
		global.Log.Error("登陆失败！用户被禁止登录")
		response.ErrorMessage(c, utils.Params, "登陆失败！用户被禁止登录")
		return
	}
	tokenNext(c, *user)
	return
}

// tokenNext 登陆后签发 JWT
func tokenNext(c *gin.Context, user model.User) {
	j := utils.NewJWT() // 唯一签名
	claims := j.CreateClaims(model.BaseClaims{
		Id:       user.Id,
		Account:  user.Account,
		Username: user.Username,
		Role:     user.Role,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.Log.Error("获取token失败!", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorTokenInvalid, "获取token失败!")
		return
	}
	// 如果没有使用多点登录，直接在当前登录的设备上设置 token
	if !global.Config.Server.UseMultipoint {
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.SuccessDetailed(c, "登录成功", response.UserLoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		})
		return
	}

	// 如果启用了多点登录
	if jwtStr, err := jwtService.GetRedisJWT(user.Account); errors.Is(err, redis.Nil) {
		if err := jwtService.SetRedisJWT(token, user.Account); err != nil {
			global.Log.Error("设置登录状态失败!", zap.Error(err))
			response.ErrorMessage(c, utils.ErrorRedis, "设置登录状态失败")
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.SuccessDetailed(c, "登录成功", response.UserLoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		})
	} else if err != nil {
		global.Log.Error("设置登录状态失败!", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorRedis, "设置登录状态失败")
		return
	} else {
		var blackJWT model.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			response.ErrorMessage(c, utils.ErrorRedis, "设置登录状态失败")
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.SuccessDetailed(c, "登录成功", response.UserLoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		})
	}
}

// GetUser
// @Tags      user
// @Summary   获取用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=model.User,msg=string}  "返回用户信息"
// @Router    /api/user/get [get]
func GetUser(c *gin.Context) {
	id := utils.GetUserId(c)
	loginUser, err := userService.GetUser(id)
	if err != nil {
		global.Log.Error("获取失败!", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorNotLogin, "获取失败")
		return
	}
	response.SuccessDetailed(c, "获取成功", loginUser)
}

// UserUpdate
// @Tags      user
// @Summary   更新用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UserUpdate                                   true  "ID,用户名,头像地址,邮箱,电话,简介,性别"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /api/user/update [post]
func UserUpdate(c *gin.Context) {
	var r request.UserUpdate
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.Error(c, utils.Params)
		return
	}
	err = userService.UserUpdate(model.User{
		Id:       r.Id,
		Username: r.Username,
		Avatar:   r.Avatar,
		Email:    r.Email,
		Phone:    r.Phone,
		Profile:  r.Profile,
		Gender:   r.Gender,
	})
	if err != nil {
		global.Log.Error("修改信息失败!", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorUpdate, "修改用户信息失败")
		return
	}
	response.SuccessMessage(c, "修改成功")
}

// UserSearch
// @Tags     user
// @Summary  分页搜索用户
// @accept   application/json
// @Produce  application/json
// @Param    data  body      request.UserSearch                                  true  "页号,页大小,排序规则,排序字段,搜索词"
// @Success  200   {object}  response.Response{data=response.PageResponse,msg=string}  "总数,每页的记录数,总页数,数据,页号"
// @Router   /api/user/search [post]
func UserSearch(c *gin.Context) {
	var r request.UserSearch
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.Error(c, utils.Params)
		return
	}
	total, size, pages, records, current, err := userService.UserSearch(r)
	if err != nil {
		global.Log.Error("获取用户列表失败", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorSearch, "获取用户列表失败")
		return
	}
	response.SuccessData(c, response.PageResponse{
		Total:   total,
		Size:    size,
		Pages:   pages,
		Records: records,
		Current: current,
	})
}

// UserDelete
// @Tags      user
// @Summary   删除用户(管理员)
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UserDelete       true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除用户"
// @Router    /api/user/delete [post]
func UserDelete(c *gin.Context) {
	var r request.UserDelete
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.Error(c, utils.Params)
		return
	}
	err = userService.UserDelete(r.Id)
	if err != nil {
		global.Log.Error("删除用户失败", zap.Error(err))
		response.ErrorMessage(c, utils.ErrorDelete, "删除用户失败")
		return
	}
	response.SuccessMessage(c, "删除成功")
}

// UploadAvatar
// @Tags      user
// @Summary   更新用户头像
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                     true  "上传文件"
// @Success   200   {object}  response.Response{msg=string}  "上传用户头像"
// @Router    /api/user/avatar [post]
func UploadAvatar(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, utils.Params)
		return
	}
	id := utils.GetUserId(c)
	// 使用使用 oss
	if global.Config.Server.UseOSS {
		if err := userService.UserAvatar(id, fileHeader); err != nil {
			response.Error(c, utils.ErrorFile)
			return
		}
	} else {
		if err := userService.UserAvatarLocal(id, fileHeader); err != nil {
			response.Error(c, utils.ErrorFile)
			return
		}
	}
	response.SuccessMessage(c, "上传成功")
}
