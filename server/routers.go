package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (svr *Server) InitRouters() *gin.Engine {

	if svr.mode == 1 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/static", "./static")
	router.Use(Cors())
	router.GET("/api/total", svr.routerGetTotal)
	//router.GET("api/hbt", svr.routerGetHbtInfo)
	router.GET("/api/pools", svr.routerGetAllPools)
	router.GET("/api/pool", svr.routerGetPoolByAddr)

	router.POST("/api/upload_file", svr.handleUploadFile)
	router.POST("/api/create_pool", svr.handlePostPoolInfo)
	router.POST("/api/modify_pool_weight", svr.handleUpdatePoolInfoWeight)
	router.GET("/api/admin/pools", svr.handleGetPoolInfo)
	//router.GET("/logs", svr.routerGetLogs)

	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-requested-with")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func Ok(c *gin.Context, datas ...interface{}) {
	if len(datas) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": datas[0],
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})

}

func Err(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"ok":   false,
		"data": data,
	})
}
