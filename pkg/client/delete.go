package client

import (
	"context"
	"net/http"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
)

func (c *UIClient) HandleDeleteShortlink(ct *gin.Context) {
	ctx, span := c.tracer.Start(ct, "ShortlinkUI.HandleNew")
	defer span.End()

	token, err := ct.Cookie(authCookieName)

	if err != nil {
		ct.AbortWithError(http.StatusUnauthorized, err)
		ct.Writer.Header().Set("Location", "/")
		ct.Writer.WriteHeader(http.StatusFound)
		return
	}

	token = strings.TrimPrefix(token, "token ")
	auth := context.WithValue(ctx, swagger.ContextAccessToken, token)

	err = ct.Request.ParseForm()
	if err != nil {
		ct.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	shortLinkName := ct.Query("name")

	client := swagger.NewAPIClient(&swagger.Configuration{
		BasePath:  "https://api.short.cedi.dev",
		UserAgent: "urlshortener-ui",
	})

	_, resp, err := client.Apiv1Api.ApiV1ShortlinkShortlinkDelete(auth, shortLinkName)
	if resp.StatusCode != http.StatusOK && err != nil {
		ct.AbortWithError(resp.StatusCode, err)
	}

	ct.Redirect(http.StatusFound, "/home")
}
