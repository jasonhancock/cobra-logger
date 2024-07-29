package logger

import (
	"io"
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
	return logger.New(
		logger.WithDestination(w),
		logger.With(keyvals...),
		logger.WithFormat(cfg.Format),
		logger.WithLevel(cfg.Level),
		logger.WithName(cfg.Name),
		logger.WithAutoCallerPrefixTrim(),
	)
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
