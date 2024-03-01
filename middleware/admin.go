package middleware

import (
	"errors"
	"github.com/axliupore/gin-template/global"
	"github.com/axliupore/gin-template/model/response"
	"github.com/axliupore/gin-template/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Admin 判断用户是不是管理员
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := utils.GetUserRole(c)
		if role != "admin" {
			global.Log.Error("权限不足", zap.Error(errors.New("权限不足")))
			response.Error(c, utils.ErrorAdmin)
			c.Abort()
			return
		}
	}
}
