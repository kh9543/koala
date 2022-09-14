package heartbeat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHandler(g *gin.RouterGroup) {
	subg := g.Group("heartbeat")
	subg.Handle("GET", "/", heartBeat)
}

func heartBeat(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}
