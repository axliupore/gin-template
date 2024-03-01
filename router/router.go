package router

import (
	"github.com/axliupore/gin-template/global"
	"github.com/axliupore/gin-template/middleware"
	"github.com/axliupore/gin-template/model/response"
	"github.com/gin-gonic/gin"
)

// Router 初始化总路由
func Router() *gin.Engine {
	if global.Config.Server.Mode == "public" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	if global.Config.Server.Mode != "public" {
		router.Use(gin.Logger())
	}

	if !global.Config.Server.UseOSS {
		router.Static(global.Config.Server.FilePath, "."+global.Config.Server.FilePath)
	}

	router.Use(middleware.Cors())

	publicGroup := router.Group(global.Config.Server.RouterPrefix)
	{
		// 健康监测
		publicGroup.GET("/health", func(c *gin.Context) {
			response.SuccessMessage(c, "ok")
		})
	}

	{
		InitUserRouter(publicGroup)
	}

	return router
}
