package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func status(c *gin.Context) {
	c.String(http.StatusOK, "All systems nominal :P")
}
