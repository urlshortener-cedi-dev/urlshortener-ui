package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func (c *UIClient) HandleNotFound(ct *gin.Context) {
	ctx := ct.Request.Context()
	span := trace.SpanFromContext(ctx)

	// Check if the span was sampled and is recording the data
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, "UIClient.HandleNotFound")
		defer span.End()
	}

	span.AddEvent("Not found")

	otelgin.HTML(
		ct,
		http.StatusNotFound,
		"404.html",
		gin.H{},
	)
}
