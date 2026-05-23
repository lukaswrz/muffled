package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/urfave/cli/v3"
	"hack.moontide.ink/pingfisher/binfo"
	"hack.moontide.ink/pingfisher/muffled/internal/events"
	"hack.moontide.ink/pingfisher/muffled/internal/listenbrainz"
	"hack.moontide.ink/pingfisher/muffled/internal/notify"
	"hack.moontide.ink/pingfisher/muffled/internal/polling"
	"hack.moontide.ink/pingfisher/muffled/internal/routing"
	"hack.moontide.ink/pingfisher/muffled/internal/widget"
)

var bi = binfo.MustGet()

const bufsize = 8

func main() {
	var p string

	cli.VersionPrinter = func(cmd *cli.Command) {
		_, _ = fmt.Fprintf(
			cmd.Root().Writer,
			"%s\n",
			bi.Summarize(
				cmd.Name,
				cmd.Version,
				binfo.Multiline|binfo.Build|binfo.VCS|binfo.Module|binfo.CGO,
			),
		)
	}

	app := &cli.Command{
		Name:    "muffled",
		Version: bi.Module.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "configuration file",
				Destination: &p,
			},
		},
		Action: func(context.Context, *cli.Command) error {
			c := configure(p)

			slog.SetLogLoggerLevel(slogLevels[c.LogLevel])

			logger := httplog.NewLogger("muffled", httplog.Options{
				LogLevel:         slogLevels[c.LogLevel],
				Concise:          true,
				RequestHeaders:   true,
				MessageFieldName: "message",
				Tags: map[string]string{
					"build":  bi.Summarize("", "", binfo.Build),
					"vcs":    bi.Summarize("", "", binfo.VCS),
					"module": bi.Summarize("", "", binfo.Module),
				},
				QuietDownRoutes: []string{"/ping"},
				QuietDownPeriod: 10 * time.Second,
				SourceFieldName: "source",
			})

			r := chi.NewRouter()
			r.Use(httplog.RequestLogger(logger))
			r.Use(middleware.Heartbeat("/ping"))

			nm := notify.NewManager[events.PlayingNowEvent](bufsize)

			widget, err := widget.Get(c.WidgetPath)
			if err != nil {
				return fmt.Errorf("unable to get widget: %w", err)
			}

			err = routing.Register(r, nm, widget)
			if err != nil {
				return fmt.Errorf("unable to register routes: %w", err)
			}

			lb, err := listenbrainz.NewClient(c.ListenBrainzBaseURL)
			if err != nil {
				log.Fatal(err)
			}
			slog.Debug("starting to poll")
			go func() {
				if err := polling.Poll(nm, lb, c.Interval, c.User); err != nil {
					fmt.Fprintf(os.Stderr, "polling aborted: %s\n", err)
					os.Exit(1)
				}
			}()

			slog.Info("starting web server", "address", c.Address)
			return http.ListenAndServe(c.Address, r)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
