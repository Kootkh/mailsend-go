package flag_and_slog

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"golang.org/x/exp/slog"
)

var loggerLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
	"fatal": slog.LevelFatal,
}

var opts struct {
	http struct {
		port   int
		secure bool
	}

	db struct {
		migrationsDir string
	}

	featureToggles struct {
		configFile string
	}

	logger struct {
		style string
		level string
	}
}

func main() {
	flag.StringVar(&opts.featureToggles.configFile, "feature-config", "featuretoggles.json", "Location of feature toggle config file")
	flag.IntVar(&opts.http.port, "http-port", 8000, "The HTTP port to run the server on")
	flag.BoolVar(&opts.http.secure, "secure", false, "Enable HTTPS on the HTTP server")
	flag.StringVar(&opts.logger.level, "logger-level", "debug", "Set logger level")
	flag.StringVar(&opts.logger.style, "logger-style", "text", "Set logger style")

	flag.Parse()

	logLevel, ok := loggerLevels[opts.logger.level]
	if !ok {
		log.Fatalf("Invalid logger level supplied: %s", opts.logger.level)
	}

	var handler slog.Handler

	switch opts.logger.style {

	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})

	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})

	case "dev":
		handler = slog.NewDevHandler(os.Stdout, logLevel)

	default:
		log.Fatalf("Invalid logger style supplied: %s", opts.logger.style)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()
	if err := cfg.Load(ctx); err != nil {
		log.Fatal("Load config: ", err)
	}

	var longTasks sync.WaitGroup

	envPort, err := strconv.ParseInt(os.Getenv("PORT"), 0, 64)
	if err != nil {

	}

}
