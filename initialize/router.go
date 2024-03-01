package initialize

import (
	"fmt"
	"github.com/axliupore/gin-template/global"
	"github.com/axliupore/gin-template/router"
)

func InitRouter() {
	routers := router.Router()
	address := fmt.Sprintf(":%d", global.Config.Server.Port)
	routers.Run(address)
}
