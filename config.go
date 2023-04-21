package logger

import (
	"io"
	"strings"

	"github.com/jasonhancock/go-env"
	"github.com/jasonhancock/go-logger"
	"github.com/spf13/cobra"
)

type Config struct {
	Level  string
	Format string
	name   string
}

func NewConfig(cmd *cobra.Command) *Config {
	pieces := strings.Fields(cmd.Use)
	c := &Config{
		name: pieces[0],
	}

	cmd.Flags().StringVar(
		&c.Level,
		"log-level",
		env.String("LOG_LEVEL", "info"),
		"Log level (all|err|warn|info|debug",
	)

	cmd.Flags().StringVar(
		&c.Format,
		"log-format",
		env.String("LOG_FORMAT", "logfmt"),
		"The format of log messages. (logfmt|json)",
	)

	return c
}

// Logger gets the logger.
func (cfg *Config) Logger(w io.Writer, keyvals ...interface{}) *logger.L {
	return logger.New(
		logger.WithDestination(w),
		logger.With(keyvals...),
		logger.WithFormat(cfg.Format),
		logger.WithLevel(cfg.Level),
		logger.WithName(cfg.name),
	)
}
