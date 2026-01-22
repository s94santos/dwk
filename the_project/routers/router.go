package routers

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes() *gin.Engine {

	r := gin.Default()

	r.GET("/", landingPage)

	return r

}

func landingPage(c *gin.Context) {
	c.File("public/index.html")
}
