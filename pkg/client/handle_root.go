package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func (c *UIClient) HandleRoot(ct *gin.Context) {
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
		log.WithError(err).Error("could not find auth cookie")

		ct.Redirect(http.StatusFound, "/login")
		return
	}

	token = strings.TrimPrefix(token, "Bearer")
	token = strings.TrimPrefix(token, "token")
	token = strings.TrimSpace(token)

	if len(token) == 0 {
		err = fmt.Errorf("could not find auth cookie")
		span.RecordError(err)
		log.WithError(err).Error("error")

		ct.Redirect(http.StatusFound, "/login")
		return
	}

	ghUser, err := getGhUser(ctx, token)
	if err != nil {
		err = fmt.Errorf("could not find auth cookie")
		span.RecordError(err)
		log.WithError(err).Error("error")

		ct.Redirect(http.StatusFound, "/login")
		return
	}

	if ghUser == nil || ghUser.Id == 0 {
		err = fmt.Errorf("could not find user on GitHub")
		span.RecordError(err)
		log.WithError(err).Error("error")

		ct.Redirect(http.StatusFound, "/login")
		return
	}

	ct.Redirect(http.StatusFound, "/home")
}
