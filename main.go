package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokryptk/metamorph/internal/cluster"
	"github.com/sokryptk/metamorph/internal/config"
	"github.com/sokryptk/metamorph/internal/tui/views"
	"github.com/spf13/pflag"
)

var (
	ctx        = context.Background()
	configPath = pflag.StringP("config", "c", "metamorph.yaml", "Path to the configuration file")
	help       = pflag.BoolP("help", "h", false, "Show this help")
)

func main() {
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	if _, err := tea.LogToFile("logfile.log", "simple"); err != nil {
		log.Fatal(err)
	}

	clusters, err := config.ImportClustersFromConfig(*configPath)
	if err != nil {
		slog.Error("error importing clusters:", slog.String("err", err.Error()))
		return
	}

	manager, err := cluster.NewManager(ctx, clusters)
	if err != nil {
		slog.Error("error creating cluster manager:", slog.String("err", err.Error()))
		return
	}

	p := tea.NewProgram(views.NewMetamorph(manager), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
