package client

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func (c *UIClient) HandleHomePage(ct *gin.Context) {
	// Extract span from the request context
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, "UIClient.HandleHomePage")
	}

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

	shortlinks, resp, err := c.apiClient.Apiv1Api.ApiV1ShortlinkGet(auth)
	if err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not query API")

		// TODO: Refactor
		if err.Error() == HttpStatusText(http.StatusUnauthorized) {
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

	otelgin.HTML(
		ct,
		http.StatusOK,
		"home.html",
		gin.H{
			"token":      token,
			"shortlinks": shortlinks,
			"copy_url":   c.config.ShortlinkURL,
		},
	)
}
