package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/cedi/urlshortener-ui/pkg/client"
)

func NewGinGonicHTTPServer(bindAddr string) (*gin.Engine, *http.Server) {
	router := gin.New()
	//router.Use(
	//otelgin.Middleware("urlshortener"),
	//secure.Secure(secure.Options{
	//	SSLRedirect:           true,
	//	SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	//	STSIncludeSubdomains:  true,
	//	FrameDeny:             true,
	//	ContentTypeNosniff:    true,
	//	BrowserXssFilter:      true,
	//	ContentSecurityPolicy: "default-src 'self' data: 'unsafe-inline'",
	//}),
	//)

	//load html file
	router.LoadHTMLGlob("html/templates/*.html")

	//static path
	router.Static("assets", "./html/assets")

	// 404 page
	router.NoRoute(func(ct *gin.Context) {
		ct.HTML(
			http.StatusNotFound,
			"404.html",
			gin.H{},
		)
	})

	srv := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}

	return router, srv
}

func Load(router *gin.Engine, uiClient *client.UIClient) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
