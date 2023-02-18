package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cedi/urlshortener-ui/pkg/config"
	"github.com/cedi/urlshortener-ui/pkg/model"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const (
	authCookieName string = "auth"
)

type UIClient struct {
	tracer trace.Tracer
	config *config.Config
}

func NewUIClient(tracer trace.Tracer, config *config.Config) *UIClient {
	return &UIClient{
		tracer: tracer,
		config: config,
	}
}

func (c *UIClient) HandleLogin(ct *gin.Context) {
	_, span := c.tracer.Start(ct, "ShortlinkUI.HandleLogin")
	defer span.End()

	redirectURI := "http://localhost:8081/oauth/redirect"
	fmt.Println(redirectURI)

	ct.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"clientID":     c.config.ClientID,
			"redirect_uri": redirectURI,
		},
	)
}

func (c *UIClient) HandleLoginOauthRedirect(ct *gin.Context) {
	ctx, span := c.tracer.Start(ct, "ShortlinkUI.HandleLoginOauthRedirect")
	defer span.End()

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{}

	// First, we need to get the value of the `code` query param
	err := ct.Request.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		ct.Writer.WriteHeader(http.StatusBadRequest)
		http.Redirect(ct.Writer, ct.Request, "/login", http.StatusTemporaryRedirect)
		return
	}
	code := ct.Request.FormValue("code")

	// Next, lets for the HTTP request to call the github oauth endpoint to get our access token
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", c.config.ClientID, c.config.ClientSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		ct.Writer.WriteHeader(http.StatusBadRequest)
		http.Redirect(ct.Writer, ct.Request, "/login", http.StatusTemporaryRedirect)
		return
	}
	// We set this header since we want the response as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		ct.Writer.WriteHeader(http.StatusInternalServerError)
		http.Redirect(ct.Writer, ct.Request, "/login", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t model.OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		ct.Writer.WriteHeader(http.StatusBadRequest)
		http.Redirect(ct.Writer, ct.Request, "/login", http.StatusTemporaryRedirect)
		return
	}

	ct.SetCookie(authCookieName, t.AccessToken, 3600, "/", "localhost", true, true)

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token
	ct.Writer.Header().Set("Location", "/home")
	ct.Writer.WriteHeader(http.StatusFound)
}
