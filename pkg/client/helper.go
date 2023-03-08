package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type githubUser struct {
	Id         int    `json:"id,omitempty"`
	Login      string `json:"login,omitempty"`
	Avatar_url string `json:"avatar_url,omitempty"`
	Type       string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
}

func getGhUser(ct context.Context, bearerToken string) (*githubUser, error) {
	span := trace.SpanFromContext(ct)

	// prepare request to GitHubs User endpoint
	req, err := http.NewRequestWithContext(ct, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Failed to build request to fetch GitHub API")
	}

	// Set headers
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+bearerToken)

	// Perform request
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Failed to fetch UserInfo from GitHub API")
	}
	defer resp.Body.Close()

	// If request was unsuccessful, we error out
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad credentials")
	}

	// If successful, we read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Error while reading the response")
	}

	// And parse it in our GithubUser model
	githubUser := &githubUser{}
	err = json.Unmarshal(body, githubUser)
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Failed to unmarshal GitHub UserInfo")
	}

	return githubUser, nil
}

func HttpStatusText(code int) string {
	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}
