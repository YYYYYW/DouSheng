package main

import (
	"DouSheng/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)
	database.Init()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// 2006-01-02 15:04

}
