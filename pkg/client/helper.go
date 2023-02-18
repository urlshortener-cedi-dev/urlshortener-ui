package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type githubUser struct {
	Id         int    `json:"id,omitempty"`
	Login      string `json:"login,omitempty"`
	Avatar_url string `json:"avatar_url,omitempty"`
	Type       string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
}

func getGhUser(c context.Context, bearerToken string) (*githubUser, error) {
	// prepare request to GitHubs User endpoint
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build request to fetch GitHub API")
	}

	// Set headers
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+bearerToken)

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
		return nil, errors.Wrap(err, "Error while reading the response")
	}

	// And parse it in our GithubUser model
	githubUser := &githubUser{}
	err = json.Unmarshal(body, githubUser)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal GitHub UserInfo")
	}

	return githubUser, nil
}

func getGhUserInfo(c context.Context, bearerToken string, username string) (*githubUser, error) {
	// prepare request to GitHubs User endpoint
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/users/"+username, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build request to fetch GitHub API")
	}

	// Set headers
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
		return nil, errors.Wrap(err, "Error while reading the response")
	}

	// And parse it in our GithubUser model
	githubUser := &githubUser{}
	err = json.Unmarshal(body, githubUser)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal GitHub UserInfo")
	}

	return githubUser, nil
}
