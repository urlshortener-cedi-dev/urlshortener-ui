package client

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
)

func (c *UIClient) HandleHomePage(ct *gin.Context) {
	ctx, span := c.tracer.Start(ct, "ShortlinkUI.HandleLogin")
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

	client := swagger.NewAPIClient(&swagger.Configuration{
		BasePath:  "https://api.short.cedi.dev",
		UserAgent: "urlshortener-ui",
	})

	shortlinks, resp, err := client.Apiv1Api.ApiV1ShortlinkGet(auth)
	if err != nil {

		// TODO: Refactor
		switch err.Error() {

		case "401 Unauthorized":
			ct.AbortWithError(http.StatusUnauthorized, err)

		default:
			ct.AbortWithError(http.StatusInternalServerError, err)
		}

		return
	}

	// TODO: Refactor
	if resp.StatusCode != 200 {
		ct.AbortWithStatus(resp.StatusCode)
		return
	}

	sort.Slice(shortlinks, func(i, j int) bool {
		return shortlinks[i].Name < shortlinks[j].Name
	})

	ct.HTML(
		http.StatusOK,
		"home.html",
		gin.H{
			"token":      token,
			"shortlinks": shortlinks,
		},
	)
}
