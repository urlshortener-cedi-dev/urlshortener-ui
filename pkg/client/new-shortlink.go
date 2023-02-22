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

func (c *UIClient) HandleNewShortlink(ct *gin.Context) {
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

	client := swagger.NewAPIClient(&swagger.Configuration{
		BasePath:  "https://api.short.cedi.dev",
		UserAgent: "urlshortener-ui",
	})

	client.Apiv1Api.ApiV1ShortlinkShortlinkPost(auth, name, swagger.V1alpha1ShortLinkSpec{
		After:  int32(redirectAfter),
		Code:   int32(code),
		Owner:  ghUser.Login,
		Owners: coOwners,
		Target: form.Get("url"),
	})

	ct.Redirect(http.StatusFound, "/home")
}
