package service

import (
	"errors"
	"github.com/axliupore/gin-template/global"
	"github.com/axliupore/gin-template/model"
	"github.com/axliupore/gin-template/model/request"
	"github.com/axliupore/gin-template/utils"
	"gorm.io/gorm"
	"math"
	"mime/multipart"
	"strconv"
	"time"
)

type UserService struct {
}

// UserRegister 用户注册
func (userService *UserService) UserRegister(user *model.User) error {
	if !errors.Is(global.Db.Where("account = ?", user.Account).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("账号已经存在")
	}
	user.Password = utils.BcryptHash(user.Password)
	err := global.Db.Create(user).Error
	return err
}

// UserLogin 用户登录
func (userService *UserService) UserLogin(user *model.User) (*model.User, error) {
	var u model.User
	if err := global.Db.Where("account = ?", user.Account).First(&u).Error; err != nil {
		return nil, errors.New("账号不存在")
	}
	if ok := utils.BcryptCheck(user.Password, u.Password); !ok {
		return nil, errors.New("密码错误")
	}
	return &u, nil
}

// GetUser 根据 id 获取用户
func (userService *UserService) GetUser(id int64) (model.User, error) {
	var user model.User
	err := global.Db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// UserUpdate 更新用户信息
func (userService *UserService) UserUpdate(user model.User) error {
	return global.Db.Model(&model.User{}).
		Select("update_at", "username", "avatar", "email", "phone", "profile", "birthday", "gender").
		Where("id = ?", user.Id).
		Updates(map[string]interface{}{
			"update_time": time.Now(),
			"username":    user.Username,
			"avatar":      user.Avatar,
			"email":       user.Email,
			"phone":       user.Phone,
			"profile":     user.Profile,
			"gender":      user.Gender,
		}).Error
}

// UserSearch 搜索用户
func (userService *UserService) UserSearch(r request.UserSearch) (int64, int, int, interface{}, int, error) {
	limit := r.PageSize
	current := r.Current
	offset := limit * (current - 1)
	sortField := r.SortField
	sortOrder := r.SortOrder
	var userList []model.User
	var total int64
	db := global.Db.Model(&model.User{})
	text := r.Text
	// 判断text是否能够转换为数字
	if _, err := strconv.ParseInt(text, 10, 64); err == nil {
		db = db.Where("id = ? OR account LIKE ? OR username LIKE ?", text, "%"+text+"%", "%"+text+"%")
	} else {
		db = db.Where("account LIKE ? OR username LIKE ?", "%"+text+"%", "%"+text+"%")
	}
	err := db.Count(&total).Error
	if err != nil {
		return 0, 0, 0, nil, 0, err
	}
	// 构建排序条件
	orderClause := sortField + " " + sortOrder
	err = db.Limit(limit).Offset(offset).Order(gorm.Expr(orderClause)).Find(&userList).Error
	// 计算总页数
	pages := int(math.Ceil(float64(total) / float64(limit)))
	return total, limit, pages, userList, current, nil
}

// UserDelete 删除用户
func (userService *UserService) UserDelete(id int64) error {
	return global.Db.Delete(&model.User{}, "id = ?", id).Error
}

// UserAvatar 更新用户头像
func (userService *UserService) UserAvatar(id int64, fileHeader *multipart.FileHeader) error {
	user, err := userService.GetUser(id)
	if err != nil {
		return err
	}

	// 上传图片到 oss
	imgURL, err := utils.UploadFile(fileHeader)
	if err != nil {
		return err
	}
	// 更新用户信息
	user.Avatar = imgURL
	if err := userService.UserUpdate(user); err != nil {
		return err
	}
	return nil
}

// UserAvatarLocal 更新用户头像到本地
func (userService *UserService) UserAvatarLocal(id int64, fileHeader *multipart.FileHeader) error {
	user, err := userService.GetUser(id)
	if err != nil {
		return err
	}
	imgURL, err := utils.SaveFileLocal(fileHeader)
	if err != nil {
		return err
	}
	user.Avatar = imgURL
	if err := userService.UserUpdate(user); err != nil {
		return err
	}
	return nil
}
