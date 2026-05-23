package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/go-multierror"
	"hack.moontide.ink/pingfisher/muffled/internal/gonfig"
)

var slogLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type config struct {
	Address             string `toml:"address"`
	LogLevel            string `toml:"log_level"`
	User                string `toml:"user"`
	Interval            int    `toml:"interval"`
	ListenBrainzBaseURL string `toml:"listenbrainz_base_url"`
	WidgetPath          string `toml:"widget_path"`
}

func unmarshal(data []byte, c *config) error {
	return toml.Unmarshal(data, &c)
}

func finalize(c *config) error {
	merr := &multierror.Error{}

	_, ok := slogLevels[c.LogLevel]
	if !ok {
		merr = multierror.Append(merr, fmt.Errorf(
			"invalid log level %q",
			c.LogLevel,
		))
	}

	if c.User == "" {
		merr = multierror.Append(merr, fmt.Errorf(
			"please specify a user",
		))
	}

	mi := 10
	if c.Interval < mi {
		merr = multierror.Append(merr, fmt.Errorf(
			"interval %d is too low, please use a value equal to or greater than %d",
			c.Interval,
			mi,
		))
	}

	return merr.ErrorOrNil()
}

func configure(path string) config {
	c := config{
		Address:             "localhost:8080",
		LogLevel:            "info",
		Interval:            120,
		ListenBrainzBaseURL: "https://api.listenbrainz.org/1",
	}

	searchPaths := []string{
		"muffled.toml",
		"/etc/muffled/muffled.toml",
	}

	_, err := gonfig.ReadConfig(path, searchPaths, &c, unmarshal, finalize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s\n", err)
		os.Exit(1)
	}

	return c
}
