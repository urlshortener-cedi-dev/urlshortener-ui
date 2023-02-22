package client

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *UIClient) HandleEditShortlink(ct *gin.Context) {
	ctx, span := c.tracer.Start(ct, "ShortlinkUI.HandleNew")
	defer span.End()

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

	client := swagger.NewAPIClient(&swagger.Configuration{
		BasePath:  "https://api.short.cedi.dev",
		UserAgent: "urlshortener-ui",
	})

	_, resp, err := client.Apiv1Api.ApiV1ShortlinkShortlinkPut(auth, name, swagger.V1alpha1ShortLinkSpec{
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
