package main

import "github.com/axliupore/gin-template/initialize"

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
func main() {
	initialize.InitConfig()
}
