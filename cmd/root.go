package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cedi/urlshortener-ui/pkg/config"
	"github.com/cedi/urlshortener-ui/pkg/observability"
	"github.com/mitchellh/go-homedir"
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
	// Version represents the Version of the kkpctl binary, should be set via ldflags -X
	Version string

	// Date represents the Date of when the kkpctl binary was build, should be set via ldflags -X
	Date string

	// Commit represents the Commit-hash from which kkpctl binary was build, should be set via ldflags -X
	Commit string

	// BuiltBy represents who build the binary, should be set via ldflags -X
	BuiltBy string

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
func Execute(version, commit, date, builtBy string) {

	// assign build flags for version info
	Version = version
	Date = date
	Commit = commit
	BuiltBy = builtBy

	var err error
	config.ConfigPath = configPath
	globalConf, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Failed to find home directory: " + err.Error())
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVar(&configPath, "config", home+"/.config/urlshortener-ui/config.yaml", "Path to the configuration file")

	rootCmd.PersistentFlags().BoolVar(&developmentMode, "develop", false, "Enable developer mode (unstructured and verbose logging)")
}
