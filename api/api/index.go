package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	app = gin.New()
	app.GET("/json", JsonHandler)
	app.GET("/today", TodayHandler)
	app.GET("/history", HistoryHandler)
	app.GET("/", JsonHandler)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
