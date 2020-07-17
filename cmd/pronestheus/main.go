package main

import (
	"fmt"
	"os"
	"pronestheus/pkg"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Version metadata set by ldflags during the build.
var (
	version string
	commit  string
	date    string
)

var cfg = &pkg.ExporterConfig{
	ListenAddr:      kingpin.Flag("listen-addr", "Address on which to expose metrics and web interface.").Default(":9999").String(),
	MetricsPath:     kingpin.Flag("metrics-path", "Path under which to expose metrics.").Default("/metrics").String(),
	Timeout:         kingpin.Flag("scrape-timeout", "The time to wait for remote APIs to response, in milliseconds.").Default("5000").Int(),
	TemperatureUnit: kingpin.Flag("temperature-unit", "The temperature metric unit [celsius, fahrenheit].").Default("celsius").Enum("celsius", "fahrenheit"),
	NestURL:         kingpin.Flag("nest-api-url", "The Nest API URL.").Default("https://developer-api.nest.com/devices/thermostats").String(),
	NestToken:       kingpin.Flag("nest-api-token", "The authorization token for Nest API.").Required().String(),
	WeatherURL:      kingpin.Flag("weather-api-url", "The OpenWeatherMap URL.").Default("http://api.openweathermap.org/data/2.5/weather").String(),
	WeatherToken:    kingpin.Flag("weather-api-token", "The authorization token for OpenWeatherMap API.").String(),
	WeatherLocation: kingpin.Flag("weather-api-location-id", "The location ID for OpenWeatherMap API. Defaults to Amsterdam.").Default("2759794").String(),
}

func main() {
	// Add short flags to --version and --help.
	kingpin.Version(versionStr()).VersionFlag.Short('v')
	kingpin.HelpFlag.Short('h')

	// Set the main command name so it can be used as a prefix for env variable names.
	kingpin.CommandLine.Name = "pronestheus"
	kingpin.CommandLine.DefaultEnvars()

	// TODO: add validators for empty values

	kingpin.Parse()

	exporter, err := pkg.NewExporter(cfg)
	exitOnErr(err)

	err = exporter.Run()
	exitOnErr(err)
}

// versionStr returns a string with version metadata: number, git sha and build date.
// It returns "development" if version variables are not set during the build.
func versionStr() string {
	if version == "" {
		return "development"
	}

	return fmt.Sprintf("%s - revision %s built at %s", version, commit[:6], date)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}