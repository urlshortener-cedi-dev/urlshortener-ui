package client

import (
	"github.com/cedi/urlshortener-ui/pkg/config"
	"github.com/cedi/urlshortener-ui/pkg/swagger"

	"go.opentelemetry.io/otel/trace"
)

const (
	authCookieName string = "auth"
)

type UIClient struct {
	tracer    trace.Tracer
	config    *config.Config
	apiClient *swagger.APIClient
}

func NewUIClient(tracer trace.Tracer, config *config.Config, apiClient *swagger.APIClient) *UIClient {
	return &UIClient{
		tracer:    tracer,
		config:    config,
		apiClient: apiClient,
	}
}
