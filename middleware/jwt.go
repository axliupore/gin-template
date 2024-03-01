package middleware

import (
	"github.com/axliupore/gin-template/model/response"
	"github.com/axliupore/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// JWT 身份认证
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 a-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := utils.GetToken(c)
		if token == "" {
			response.ErrorMessage(c, utils.ErrorTokenAuth, "未登录或非法访问")
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, code := j.ParseToken(token)
		// 解析错误
		if code != utils.Success {
			response.Error(c, code)
			utils.ClearToken(c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
		// todo 这里可以对即将过期的 token 进行刷新操作
	}
}
