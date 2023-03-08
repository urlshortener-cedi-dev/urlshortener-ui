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

func (c *UIClient) HandleEdit(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, "UIClient.HandleEdit")
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

	shortLink, resp, err := c.apiClient.Apiv1Api.ApiV1ShortlinkShortlinkGet(auth, shortLinkName)
	if resp.StatusCode != http.StatusOK && err != nil {
		ct.AbortWithError(resp.StatusCode, err)
	}

	otelgin.HTML(
		ct,
		http.StatusOK,
		"edit.html",
		gin.H{
			"edit_mode": true,
			"shortlink": shortLink,
			"coowners":  strings.Join(shortLink.Spec.Owners, ", "),
		},
	)
}

func (c *UIClient) HandleEditShortlink(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, "UIClient.HandleEditShortlink")
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
	owner := form.Get("owner")
	coOwners := strings.Split(form.Get("co-owners"), ",")
	redirectType := form.Get("redirectTypeOption")
	redirectAfter, _ := strconv.Atoi(form.Get("redirectAfter"))
	code, _ := strconv.Atoi(form.Get("httpStatusCode"))

	if redirectType == "html" {
		code = 200
	}

	_, resp, err := c.apiClient.Apiv1Api.ApiV1ShortlinkShortlinkPut(auth, name, swagger.V1alpha1ShortLinkSpec{
		After:  int32(redirectAfter),
		Code:   int32(code),
		Owner:  owner,
		Owners: coOwners,
		Target: form.Get("url"),
	})
	if resp.StatusCode != 200 && err != nil {
		ct.AbortWithError(resp.StatusCode, err)
	}

	ct.Redirect(http.StatusFound, "/home")
}
