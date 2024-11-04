package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	handler "caibai/api"
)

func main() {
	// 创建 gin 引擎
	app := gin.Default()

	// 注册路由
	app.GET("/json", handler.JsonHandler)
	app.GET("/today", handler.TodayHandler)
	app.GET("/history", handler.HistoryHandler)
	app.GET("/", handler.JsonHandler)

	// 设置服务器
	srv := &http.Server{
		Addr:    ":8080", // 可以通过环境变量配置端口
		Handler: app,
	}

	// 优雅关闭
	go func() {
		// 监听系统信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		// 创建一个5秒超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
	}()

	// 启动服务器
	log.Printf("Server is running on http://localhost%s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
