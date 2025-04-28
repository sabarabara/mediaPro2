//This file is routing apiendpoint
package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)
func SetupRouter() *gin.Engine {
	r := gin.Default()

	//testroute
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "API is up and running!")
	})
	//apiroute
	return r
}
