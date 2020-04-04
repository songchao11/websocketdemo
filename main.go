package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/api/websocket", WsHandler)
	r.Run(":8080")
}
