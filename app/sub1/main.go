/*
@author '彼时思默'
@time 2020/5/11 下午5:51
@describe:
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"sync"
)

func main() {
	g:=gin.Default()
	g.POST("/post", func(context *gin.Context) {
		var msg map[string]interface{}
		err:=context.BindJSON(&msg)
		if err != nil {
			logrus.Error("json error:",err)
		}
		fmt.Println("msg:",msg)
	})
	go func() {
		err:=g.Run(":8888")
		if err != nil {
			logrus.Error("run error:",err)
		}
	}()

	wg:=sync.WaitGroup{}
	wg.Add(1)
	body:=map[string]interface{}{
		"id":"http1",
		"acceptType":"http",
		"isFlow":true,
		"path":"http://localhost:8888/post",
	}
	client :=resty.New()
	_,err:=client.NewRequest().SetBody(&body).Post("http://localhost:8080/api/register/sub")
	if err != nil {
		logrus.Error("request error:",err)
	}
	wg.Wait()
}