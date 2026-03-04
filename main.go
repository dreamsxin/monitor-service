package main

import (
	"log"
	"monitor-service/db"
	"monitor-service/handlers"
	"net/http"

	_ "monitor-service/docs" // 导入 swagger 文档

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           站点监控服务 API
// @version         1.0
// @description     用于接收监控工具上报的变化数据，并提供站点 Hook 脚本配置。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

func main() {
	// 初始化 SQLite 数据库
	db.InitDB("monitor.db")

	r := gin.Default()

	// 静态文件和模板
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/upload", handlers.UploadChange)
		api.GET("/sites", handlers.ListSites)
		api.GET("/changes", handlers.GetChanges)
		api.GET("/content", handlers.GetFileContent)
		api.GET("/site-hooks", handlers.GetSiteHooks)
	}

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Println("Server started at :8080")
	log.Println("Swagger docs at http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}
