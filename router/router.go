package router

import (
	"github.com/axliupore/gin-template/global"
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

	return router
}
