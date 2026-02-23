package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kexirong/coredns-admin/controller"
	"github.com/kexirong/coredns-admin/middleware"
)

func initRoute() {
	Router.StaticFile("/", "./dist/index.html")
	Router.Static("/js", "./dist/js")
	Router.Static("/css", "./dist/css")
	Router.Static("/fonts", "./dist/fonts")
	Router.StaticFile("/favicon.ico", "./dist/favicon.ico")
	Router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "GET" {
			c.File("./dist/index.html")
		}
	})
	Router.Use(middleware.Cors())
	Router.POST("/login", controller.Login)

	var v1 = Router.Group("/api/v1", middleware.JWTAuth())

	v1.GET("/records", controller.GetRecords)
	v1.POST("/record", controller.PostRecord)
	v1.DELETE("/record/:key", controller.DeleteRecord)
	v1.PUT("/record/:key", controller.PutRecord)

	v1.GET("/domains", controller.GetDomains)
}
