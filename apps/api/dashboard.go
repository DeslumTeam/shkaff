package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{})
}
