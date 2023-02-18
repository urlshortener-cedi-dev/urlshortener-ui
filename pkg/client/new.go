package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *UIClient) HandleNew(ct *gin.Context) {
	_, span := c.tracer.Start(ct, "ShortlinkUI.HandleNew")
	defer span.End()

	ct.HTML(
		http.StatusOK,
		"new.html",
		gin.H{},
	)
}
