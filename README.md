<div align="center">
    <img src='./docs/log.png' style='width: 15%'/>
</div>

<h1 align="center">gin-template</h1>

## 🚀技术栈

|     技术     |      功能      |
| :----------: | :------------: |
|     gin      |    web 框架    |
|     gorm     | 对象关系映射库 |
|    mysql     |     数据库     |
|    redis     |   缓存数据库   |
| aliyun - oss | 阿里云对象存储 |
|     zap      |     日志库     |
|   swagger    |  自动生成文档  |

## 👨🏻‍💻功能

- 用户管理
- 对象存储
- 文件管理
- 基于令牌的鉴权
- 集成了日志库
- 多点登录
- 缓存

## 🎄特点

-   📦 开箱即用，只需修改配置文件
-   📝 全面注释说明，学习低成本
-   🚀 启动编译迅速
-   🌱 极易定制, 拓展容易

## 📂项目结构

```
├── gin-template
    ├── api           [api 接口层]
    ├── config        [配置文件]
    ├── core          [核心文件]
    ├── docs          [swagger 文档目录]
    ├── global        [全局对象]  
    ├── initialize    [初始化]
    ├── middleware    [中间件]
    ├── model         [实体对象层]
    │   ├── request   [请求结构体]
    │   └── response  [响应结构体]                  
    ├── router        [路由层]
    ├── service       [service 层]
    ├── sql           [数据库]
    └── utils         [工具包]
```

## ⌛️快速开始

```bash
# 克隆项目
git@github.com:AxLiupore/gin-template.git

# 进入文件夹
cd gin-template

# 安装依赖
go generate

# 编译
go build -o server main.go
```

