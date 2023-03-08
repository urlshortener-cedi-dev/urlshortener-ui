package client

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func (c *UIClient) HandleNew(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, "UIClient.HandleNew")
		defer span.End()
	}

	otelgin.HTML(
		ct,
		http.StatusOK,
		"new.html",
		gin.H{},
	)
}

func (c *UIClient) HandleNewShortlink(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, "UIClient.HandleNewShortlink")
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

	form := ct.Request.Form

	name := form.Get("name")
	coOwners := strings.Split(form.Get("co-owners"), ",")
	redirectType := form.Get("redirectTypeOption")
	redirectAfter, _ := strconv.Atoi(form.Get("redirectAfter"))
	code, _ := strconv.Atoi(form.Get("httpStatusCode"))

	ghUser, err := getGhUser(ctx, token)
	if err != nil {
		ct.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if redirectType == "html" {
		code = 200
	}

	c.apiClient.Apiv1Api.ApiV1ShortlinkShortlinkPost(auth, name, swagger.V1alpha1ShortLinkSpec{
		After:  int32(redirectAfter),
		Code:   int32(code),
		Owner:  ghUser.Login,
		Owners: coOwners,
		Target: form.Get("url"),
	})

	ct.Redirect(http.StatusFound, "/home")
}
