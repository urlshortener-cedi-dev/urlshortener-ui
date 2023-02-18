package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cedi/urlshortener-ui/pkg/config"
	"github.com/cedi/urlshortener-ui/pkg/observability"
	"github.com/spf13/cobra"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	log "github.com/sirupsen/logrus"
)

const (
	serviceName    = "urlshortener-ui"
	serviceVersion = "0.0.1"
)

var (
	globalConf      *config.Config
	configPath      string
	developmentMode bool
	Tracer          trace.Tracer
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "urlshortener-ui",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	config.ConfigPath = configPath

	globalConf = config.NewConfig()

	globalConf.ClientID = os.Getenv("CLIENT_ID")
	globalConf.ClientSecret = os.Getenv("CLIENT_SECRET")
	globalConf.SessionState = os.Getenv("SESSION_STATE")
	globalConf.RedirectURL = os.Getenv("REDIRECT_URL")
	globalConf.HostName = os.Getenv("HOSTNAME")

	if developmentMode {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: false,
		})
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}

	// Initialize Tracing (OpenTelemetry)
	var tp *sdkTrace.TracerProvider
	tp, Tracer, err = observability.InitTracer(serviceName, serviceVersion)
	if err != nil {
		log.Error(err, "failed initializing tracing")
		os.Exit(1)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error(err, "Error shutting down tracer provider")
		}
	}()

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&developmentMode, "develop", false, "Enable developer mode (unstructured and verbose logging)")
}
