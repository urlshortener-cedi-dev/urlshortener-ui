package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"github.com/cedi/urlshortener-ui/pkg/client"
	promRouter "github.com/cedi/urlshortener/pkg/router"
)

func NewGinGonicHTTPServer(bindAddr string, tracer trace.Tracer, serviceName string) (*gin.Engine, *http.Server) {
	router := gin.New()
	router.Use(
		otelgin.Middleware(serviceName),
		promRouter.PromMiddleware(serviceName),
		//secure.Secure(secure.Options{
		//	SSLRedirect:           true,
		//	SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		//	STSIncludeSubdomains:  true,
		//	FrameDeny:             true,
		//	ContentTypeNosniff:    true,
		//	BrowserXssFilter:      true,
		//	ContentSecurityPolicy: "default-src 'self' data: 'unsafe-inline'",
		//}),
	)

	//load html file
	router.LoadHTMLGlob("html/templates/*.html")

	//static path
	router.Static("assets", "./html/assets")

	srv := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}

	return router, srv
}

func Load(router *gin.Engine, uiClient *client.UIClient) {
	router.NoRoute(uiClient.HandleNotFound) // 404 page

	router.GET("/", uiClient.HandleLogin)
	router.GET("/login", uiClient.HandleLogin)
	router.GET("/oauth/redirect", uiClient.HandleLoginOauthRedirect)

	router.GET("/home", uiClient.HandleHomePage)

	router.GET("/new", uiClient.HandleNew)
	router.GET("/edit", uiClient.HandleEdit)
	router.GET("/delete", uiClient.HandleDeleteShortlink)

	router.POST("/new/shortlink", uiClient.HandleNewShortlink)
	router.POST("/edit/shortlink", uiClient.HandleEditShortlink)

}
