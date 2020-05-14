/*
@author '彼时思默'
@time 2020/4/29 上午9:47
@describe:
*/
package proxy

import (
	"fmt"
	. "github.com/bishisimo/omc_log_system/core"
	"github.com/bishisimo/omc_log_system/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateProxy(port uint) {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.Cors())
	router.Use(middleware.Logger2File())
	creatRouter(router)
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		println(err)
	}
}

func creatRouter(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		register := api.Group("/register")
		{
			register.POST("/pub/:id",
				func(context *gin.Context) {
					Redux.AddPub(context.Param("id"))
					context.String(http.StatusOK, "ok")
					return
				})
			register.POST("/sub",
				func(context *gin.Context) {
					so := NewSubOption()
					err := context.BindJSON(&so)
					if err != nil {
						context.String(http.StatusBadRequest, "Json format error")
						return
					}
					if !so.IsInfoRight() {
						context.String(http.StatusBadRequest, "info miss")
						return
					}
					Redux.AddSub(so)
					context.String(http.StatusOK, "ok")
					return
				})
		}
		deregister := api.Group("/deregister")
		{
			deregister.POST("/pub/:id",
				func(context *gin.Context) {
					Redux.RemovePub(context.Param("id"))
					context.String(http.StatusOK, "ok")
					return
				})
			deregister.POST("/sub/:id",
				func(context *gin.Context) {
					Redux.RemoveSub(context.Param("id"))
					context.String(http.StatusOK, "ok")
					return
				})
		}
		pub := api.Group("/pub")
		{
			pub.POST("/:id", func(context *gin.Context) {
				msg := make(map[string]interface{})
				if pub := Redux.Pubs.Load(context.Param("id")); pub != nil {
					err := context.BindJSON(&msg)
					if err != nil {
						context.String(http.StatusBadRequest, "Json format error")
						return
					}
					pub.MsgChan <- msg
					context.String(http.StatusOK, "ok")
					return
				} else {
					context.String(http.StatusUnauthorized, "unregister")
					return
				}
			})
		}
		show := api.Group("/show")
		{
			show.GET("", func(context *gin.Context) {
				result := make(map[string]interface{})

				pubInfo := make([]string, 0)
				Redux.Pubs.Range(func(id string, value *Pub) {
					pubInfo = append(pubInfo, id)
				})
				flowInfo := make(map[string]interface{})
				Redux.SubsFlow.Range(func(id string, sub *Sub) {
					flowInfo[id] = sub.GetInfo()
				})
				batchInfo := make(map[string]interface{})
				Redux.SubsBatch.Range(func(id string, sub *Sub) {
					batchInfo[id] = sub.GetInfo()
				})
				result["pubInfo"] = pubInfo
				result["flowInfo"] = flowInfo
				result["batchInfo"] = batchInfo
				result["supportInfo"] = SubSupport
				context.JSON(http.StatusOK, result)
				return
			})
		}
	}
}
