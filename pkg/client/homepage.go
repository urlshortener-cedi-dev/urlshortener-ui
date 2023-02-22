package client

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *UIClient) HandleHomePage(ct *gin.Context) {
	ctx, span := c.tracer.Start(ct, "ShortlinkUI.HandleLogin")
	defer span.End()

	log := logrus.WithContext(ctx)

	token, err := ct.Cookie(authCookieName)
	if err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not parse auth cookie")

		ct.Redirect(http.StatusFound, "/login")
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
		span.RecordError(err)
		log.WithError(err).Error("could not query API")

		// TODO: Refactor
		if err.Error() == "401 Unauthorized" {
			ct.AbortWithStatus(http.StatusUnauthorized)
		} else {
			ct.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	// TODO: Refactor
	if resp.StatusCode != 200 {
		ct.AbortWithStatus(resp.StatusCode)
		log.WithField("StatusCode", resp.StatusCode).Error("API request was not successful")
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
			"copy_url":   c.config.ShortlinkURL,
		},
	)
}
