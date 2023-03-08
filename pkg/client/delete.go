package client

import (
	"context"
	"net/http"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func (c *UIClient) HandleDeleteShortlink(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, "UIClient.HandleDeleteShortlink")
		defer span.End()
	}

	log := logrus.New().WithContext(ctx)

	token, err := ct.Cookie(authCookieName)
	if err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not parse auth cookie")

		ct.Redirect(http.StatusFound, "/login")
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

	_, resp, err := c.apiClient.Apiv1Api.ApiV1ShortlinkShortlinkDelete(auth, shortLinkName)
	if resp.StatusCode != http.StatusOK && err != nil {
		ct.AbortWithError(resp.StatusCode, err)
	}

	ct.Redirect(http.StatusFound, "/home")
}
