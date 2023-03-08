package cmd

import (
	"net/http"

	"github.com/cedi/urlshortener-ui/pkg/client"
	"github.com/cedi/urlshortener-ui/pkg/router"
	"github.com/cedi/urlshortener-ui/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	log "github.com/sirupsen/logrus"
)

var (
	bindAddress string
	debug       bool
)

var serveCMD = &cobra.Command{
	Use:     "serve",
	Short:   "Serve the Webserver",
	Example: "urlshortener-ui serve",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		apiClient := swagger.NewAPIClient(&swagger.Configuration{
			BasePath:  globalConf.ShortlinkURL,
			UserAgent: "urlshortener-ui",
			HTTPClient: &http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			},
		})

		uiClient := client.NewUIClient(Tracer, globalConf, apiClient)

		// Init Gin Framework
		if debug {
			gin.SetMode(gin.DebugMode)

			log.SetLevel(log.DebugLevel)
			log.SetFormatter(&log.TextFormatter{})
		} else {
			gin.SetMode(gin.ReleaseMode)

			log.SetLevel(log.InfoLevel)
			log.SetFormatter(&log.JSONFormatter{})
		}

		r, srv := router.NewGinGonicHTTPServer(bindAddress, Tracer)

		router.Load(r, uiClient)

		if err := srv.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCMD)
	serveCMD.Flags().StringVarP(&bindAddress, "bind-address", "p", ":8080", "Bind Address")
	serveCMD.Flags().BoolVarP(&debug, "debug", "d", false, "Enable Debug Mode")
}
