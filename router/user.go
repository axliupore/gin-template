package router

import (
	"github.com/axliupore/gin-template/api"
	"github.com/axliupore/gin-template/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	{
		userRouter.POST("register", api.UserRegister)
		userRouter.POST("login", api.UserLogin)
		userRouter.POST("search", api.UserSearch)

		userRouter.Use(middleware.JWT())
		userRouter.GET("get", api.GetUser)
		userRouter.POST("update", api.UserUpdate)
		userRouter.POST("avatar", api.UploadAvatar)

		userRouter.Use(middleware.Admin())
		userRouter.POST("delete", api.UserDelete)
	}
}
