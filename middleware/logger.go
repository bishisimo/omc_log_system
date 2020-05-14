/*
@author '彼时思默'
@time 2020/4/2 13:42
@describe:
*/
package middleware

import (
	"github.com/gin-gonic/gin"
	rotateLogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func Logger2File() gin.HandlerFunc {
	logger := logrus.New()
	src, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logDir := "logs/"
	if _, err := os.Stat(logDir); !os.IsExist(err) {
		_ = os.Mkdir(logDir, os.ModePerm)
	}
	appName := "LogSystem"
	logPath := path.Join(logDir, appName+"%Y-%m-%d-%H-%M.log")
	logWriter, _ := rotateLogs.New(
		logPath,
		rotateLogs.WithLinkName(appName), // 生成软链，指向最新日志文件
		rotateLogs.WithMaxAge(30*24*time.Hour),    // 文件最大保存时间
		rotateLogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)
		urlPath := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logger.WithFields(
			logrus.Fields{
				"statusCode": statusCode,
				"latency":    latency,
				"clientIP":   clientIP,
				"method":     method,
				"path":       urlPath,
			}).Info()
	}
}
func Logger2Mongo() {

}
func Logger2ES() {

}
func Logger2MQ() {

}
