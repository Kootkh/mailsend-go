package flag_and_slog

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type DevHandler struct {
	level slog.Leveler
	group string
	attrs []slog.Attr
	mu    *sync.Mutex
	w     io.Writer
}

func NewDevHandler(w io.Writer, level slog.Leveler) *DevHandler {

	if level == nil || reflect.TypeOf(level).Kind() == reflect.Ptr && reflect.ValueOf(level).IsNil() {
		level = &slog.LevelVar{}
	}

	return &DevHandler{
		level: level,
		mu:    new(sync.Mutex{}),
		w:     w,
	}

}

func (h *DevHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *DevHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var attrs string

	for _, a := range h.attrs {
		if !a.Equal(slog.Attr{}) {
			attrs += " "

			if h.group != "" {
				attrs += h.group + "."
			}

			attrs += a.Key + ": " + a.Value.String() + "\n"
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		if !a.Equal(slog.Attr{}) {
			attrs += " "

			if h.group != "" {
				attrs += h.group + "."
			}

			attrs += a.Key + ": " + a.Value.String() + "\n"
		}

		return true
	})

	attrs = strings.TrimRight(attrs, "\n")

	var newlines string

	if attrs != "" {
		newlines = "\n\n"
	}

	fmt.Fprintf(h.w, "[%v] %v\n%v%v", r.Time.Format("15:04:05 MST"), r.Message, attrs, newlines)

	return nil
}

func (h *DevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DevHandler{
		level: h.level,
		group: h.group,
		attrs: append(h.attrs, attrs...),
		mu:    h.mu,
		w:     h.w,
	}
}

func (h *DevHandler) WithGroup(name string) slog.Handler {
	return &DevHandler{
		level: h.level,
		group: strings.TrimSuffix(name+"."+h.group, "."),
		attrs: h.attrs,
		mu:    h.mu,
		w:     h.w,
	}
}

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
