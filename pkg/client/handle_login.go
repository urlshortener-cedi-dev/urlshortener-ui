package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cedi/urlshortener-ui/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *UIClient) HandleLogin(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, "UIClient.HandleLogin")
		defer span.End()
	}

	redirectURI := c.config.RedirectURL
	span.SetAttributes(attribute.String("redirect_uri", redirectURI))

	otelgin.HTML(
		ct,
		http.StatusOK,
		"login.html",
		gin.H{
			"clientID":     c.config.ClientID,
			"redirect_uri": redirectURI,
		},
	)
}

func (c *UIClient) HandleLoginOauthRedirect(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, "UIClient.HandleLoginOauthRedirect")
		defer span.End()
	}

	log := logrus.WithContext(ctx)

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// First, we need to get the value of the `code` query param
	err := ct.Request.ParseForm()
	if err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not parse query")

		ct.Redirect(http.StatusFound, "/login")
		return
	}
	code := ct.Request.FormValue("code")

	if len(code) == 0 {
		errMsg := ct.Request.FormValue("error")
		errDesc := ct.Request.FormValue("error_description")
		errURI := ct.Request.FormValue("error_uri")

		err = fmt.Errorf("%s: %s", errMsg, errDesc)

		span.RecordError(err)
		log.WithError(err).WithFields(logrus.Fields{
			"error":             errMsg,
			"error_description": errDesc,
			"error_uri":         errURI,
		}).Error("Failed GitHub Auth request")

		ct.Redirect(http.StatusFound, "/login")
		return
	}

	// Next, lets for the HTTP request to call the github oauth endpoint to get our access token
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", c.config.ClientID, c.config.ClientSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not create HTTP request")

		ct.Redirect(http.StatusFound, "/login")
		return
	}
	// We set this header since we want the response as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not send HTTP request")

		ct.Redirect(http.StatusFound, "/login")
		return
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t model.OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		span.RecordError(err)
		log.WithError(err).Error("could not parse JSON response")

		ct.Redirect(http.StatusFound, "/login")
		return
	}

	if len(t.AccessToken) == 0 {
		otelgin.HTML(
			ct,
			http.StatusInternalServerError,
			"500.html",
			gin.H{},
		)
	}

	ct.SetCookie(authCookieName, t.AccessToken, 3600, "/", c.config.DashboardURL, true, true)
	ct.Redirect(http.StatusFound, "/home")
}
