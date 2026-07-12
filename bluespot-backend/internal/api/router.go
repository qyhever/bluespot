package router

import (
	"fmt"
	"net/http"
	"runtime"

	_ "bluespot/docs"
	"bluespot/internal/config"
	"bluespot/internal/controller"
	"bluespot/internal/middleware"
	"bluespot/internal/repository/persistence"
	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
	knife4goFiles "github.com/go-webtools/knife4go"
	knife4goGin "github.com/go-webtools/knife4go/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	isProd := config.IsProduction()
	// Gin 开启生产模式(默认是debug模式，会输出大量调试日志)
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 静态文件服务
	r.Static("/public", "./public")
	// /api/swagger/index.html
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// /api/k4/index.html
	r.GET("/api/k4/*any", knife4goGin.WrapHandler(knife4goFiles.Handler))

	fmt.Printf("Go Version %v\n", runtime.Version())

	authService := service.NewAuthService()
	authController := controller.NewAuthController(authService)

	metaController := controller.NewMetaController()
	appRepo := persistence.NewAppRepository()
	appService := service.NewAppService(appRepo)
	appController := controller.NewAppController(appService)
	userRepo := persistence.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	attachRepo := persistence.NewAttachRepository()
	attachService := service.NewAttachService(attachRepo)
	attachController := controller.NewAttachController(attachService)
	uploadRepo := persistence.NewUploadRepository()
	uploadService := service.NewUploadService(uploadRepo)
	uploadController := controller.NewUploadController(uploadService)
	mailService := service.NewMailService()
	mailController := controller.NewMailController(mailService)

	v1 := r.Group("/api")

	v1.GET("/meta", metaController.GetMeta)
	v1.POST("/mail", mailController.Send)

	appGroup := v1.Group("/app")
	{
		appGroup.POST("/getHelloInfo", appController.GetHelloInfo)
	}

	authGroup := v1.Group("/auth")
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/refresh", authController.RefreshToken)

	userGroup := v1.Group("/user")
	userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.GET("/info", userController.GetCurrentUserInfo)

	attachGroup := v1.Group("/attach")
	attachGroup.Use(middleware.JWTAuthMiddleware())
	attachGroup.POST("/upload", attachController.Upload)

	uploadGroup := v1.Group("/upload")
	uploadGroup.Use(middleware.JWTAuthMiddleware())
	uploadGroup.POST("/verify", uploadController.Verify)
	uploadGroup.POST("/chunk", uploadController.UploadChunk)
	uploadGroup.POST("/merge", uploadController.Merge)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}
