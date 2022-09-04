package logger

import (
	"io"
	"strings"

	"github.com/jasonhancock/go-env"
	"github.com/jasonhancock/go-logger"
	"github.com/spf13/cobra"
)

type Config struct {
	Level string
	name  string
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

	return c
}

// Logger gets the logger.
func (cfg *Config) Logger(w io.Writer, keyvals ...interface{}) *logger.L {
	return logger.New(w, cfg.name, cfg.Level, keyvals...)
}
