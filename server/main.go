package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"x_admin/config"
	"x_admin/core"
	"x_admin/core/response"
	"x_admin/middleware"
	"x_admin/routers"

	_ "x_admin/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFs embed.FS

// initRouter 初始化router
func initRouter() *gin.Engine {
	// 初始化gin
	gin.SetMode(config.Config.GinMode)
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// 设置静态路径
	router.Static(config.Config.PublicPrefix, config.Config.UploadDirectory)

	staticHttpFs := http.FS(staticFs)
	router.GET("/api/static/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		fmt.Println(filepath)

		c.FileFromFS("static"+filepath, staticHttpFs)
	})
	// 设置中间件
	router.Use(gin.Logger(), middleware.Cors(), middleware.ErrorRecover())
	router.GET("/api/admin/apiList", middleware.TokenAuth(), func(ctx *gin.Context) {
		var path = []string{}
		for _, route := range router.Routes() {
			// fmt.Printf("%s 127.0.0.1:%v%s\n", route.Method, config.Config.ServerPort, route.Path)
			path = append(path, route.Path)
		}
		response.Result(ctx, response.Success, path)
	})

	// 演示模式
	if config.Config.DisallowModify {
		router.Use(middleware.ShowMode())
	}
	// 特殊异常处理
	router.NoMethod(response.NoMethod)
	// router.NoRoute(response.NoRoute)
	// 注册路由
	apiGroup := router.Group("/api")

	routers.RegisterGroup(apiGroup)

	return router
}

// initServer 初始化server
func initServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           ":" + strconv.Itoa(config.Config.ServerPort),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

//	@title			x_admin文档
//	@version		0.0.1
//	@description	x_admin是一个完整的后台管理系统
//	@termsOfService	http://x.adtk.cn

//	@contact.name	API Support
//	@contact.url	http://x.adtk.cn
//	@contact.email	11675084@qq.com

//	@license.name	MIT License
//	@license.url	https://gitee.com/xiangheng/x_admin/blob/main/LICENSE

//	@host		localhost:8001
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// 刷新日志缓冲
	defer core.Logger.Sync()
	// 程序结束前关闭数据库连接
	if core.GetDB() != nil {
		db, _ := core.GetDB().DB()
		defer db.Close()
	}

	// 初始化router
	router := initRouter()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	fmt.Println("格式化文档注释:", "swag fmt")
	fmt.Println("生成文档:", "swag init")

	fmt.Printf("文档: http://localhost:%v/swagger/index.html", config.Config.ServerPort)
	// 初始化server
	s := initServer(router)
	// 运行服务
	log.Fatalln(s.ListenAndServe().Error())
}
