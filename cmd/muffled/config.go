package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/go-multierror"
	"hack.moontide.ink/lukas/gonfig"
)

var slogLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type config struct {
	Listen              string `toml:"listen"`
	LogLevel            string `toml:"log_level"`
	User                string `toml:"user"`
	Interval            int    `toml:"interval"`
	ListenBrainzBaseURL string `toml:"listenbrainz_base_url"`
	WidgetPath          string `toml:"widget_path"`
}

func (c *config) Unmarshal(data []byte) error {
	return toml.Unmarshal(data, c)
}

func (c *config) Validate() error {
	merr := &multierror.Error{}

	if c.Listen == "" {
		merr = multierror.Append(merr, errors.New(
			"the listen address cannot be empty",
		))
	}

	_, ok := slogLevels[c.LogLevel]
	if !ok {
		merr = multierror.Append(merr, fmt.Errorf(
			"invalid log level %q",
			c.LogLevel,
		))
	}

	if c.User == "" {
		merr = multierror.Append(merr, errors.New(
			"please specify a user",
		))
	}

	mi := 1
	if c.Interval < mi {
		merr = multierror.Append(merr, fmt.Errorf(
			"interval %d is too low, please use a value equal to or greater than %d",
			c.Interval,
			mi,
		))
	}

	if c.ListenBrainzBaseURL == "" {
		merr = multierror.Append(merr, errors.New(
			"the ListenBrainz base URL cannot be empty",
		))
	}

	return merr.ErrorOrNil()
}

func configure(exp string) config {
	c := config{
		Listen:              "localhost:8080",
		LogLevel:            "info",
		Interval:            120,
		ListenBrainzBaseURL: "https://api.listenbrainz.org/1",
	}

	search := []string{
		".",
		"/etc/muffled",
	}

	_, err := gonfig.ReadConfig(&c, "muffled.toml", search, exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s\n", err)
		os.Exit(1)
	}

	return c
}
