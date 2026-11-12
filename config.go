package logger

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/jasonhancock/go-env"
	"github.com/jasonhancock/go-helpers"
	"github.com/jasonhancock/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Config struct {
	Level  string
	Format string
	Name   string

	logger     *logger.L
	logLeveler LevelSetter
}

func NewConfig(cmd *cobra.Command) *Config {
	return NewConfigPflags(strings.Fields(cmd.Use)[0], cmd.Flags())
}

func NewConfigPflags(appName string, flags *pflag.FlagSet) *Config {
	c := Config{Name: appName}

	const envLogLevel = "LOG_LEVEL"
	flags.StringVar(
		&c.Level,
		"log-level",
		env.String(envLogLevel, "info"),
		helpers.EnvDesc("Log level (all|err|warn|info|debug).", envLogLevel),
	)

	const envLogFormat = "LOG_FORMAT"
	flags.StringVar(
		&c.Format,
		"log-format",
		env.String(envLogFormat, logger.FormatLogFmt),
		helpers.EnvDesc("The format of log messages ("+strings.Join(logger.AvailableFormats, "|")+").", envLogFormat),
	)

	return &c
}

// Logger gets the logger.
func (cfg *Config) Logger(w io.Writer, keyvals ...interface{}) *logger.L {
	return cfg.LoggerWithOptions(w, logger.With(keyvals...))
}

// LoggerWithOptions gets the logger, passing opts through to logger.New after
// the options derived from the config, so they can add to them (for example,
// logger.WithContextResolutionFunc) or override them. Like Logger, the logger
// is constructed on the first call and cached; opts passed to later calls are
// ignored.
func (cfg *Config) LoggerWithOptions(w io.Writer, opts ...logger.Option) *logger.L {
	if cfg.logger == nil {
		logLeveler := logger.NewDynamicLeveler(cfg.Level)
		cfg.logLeveler = logLeveler
		cfg.logger = logger.New(append([]logger.Option{
			logger.WithDestination(w),
			logger.WithFormat(cfg.Format),
			logger.WithName(cfg.Name),
			logger.WithAutoCallerPrefixTrim(),
			logger.WithLeveler(logLeveler),
		}, opts...)...)
	}

	return cfg.logger
}

// GetLoggerName traverses cobra commands and builds a period delimited string
// useful for setting the logger name.
func GetLoggerName(cmd *cobra.Command) string {
	return strings.Join(getCmdPath(cmd), ".")
}

func getCmdPath(cmd *cobra.Command) []string {
	name := strings.Fields(cmd.Use)[0]
	result := []string{name}
	if cmd.HasParent() {
		result = append(getCmdPath(cmd.Parent()), name)
	}

	return result
}

type LevelSetter interface {
	SetLevel(level string)
}

type LogLevelChangeRequest struct {
	Level string `json:"level"`
}

func (cfg *Config) LogLevelHandler(w http.ResponseWriter, r *http.Request) {
	var req LogLevelChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cfg.logLeveler.SetLevel(req.Level)
	w.WriteHeader(http.StatusNoContent)
}
