package main

import (
	"fmt"
	"log/slog"
	"os"
	"redis-challenge/internal/config"
	"redis-challenge/internal/server"
	"redis-challenge/internal/store"
)

func main() {
	configuration, err := config.LoadConfiguration()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		os.Exit(1)
	}

	serverMonitor := make(server.MonitorChannel)

	srv, err := server.NewChallengeServer(configuration.Port, store.NewBuilder()).
		RestoreFromArchive(configuration.AppendLogReader).
		WithArchiveWriter(configuration.AppendLogWriter).
		WithMonitorChannel(serverMonitor).
		Start()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create server: %v", err))
		os.Exit(1)
	}

	slog.Info("Listening on", "address", srv.Address())

	for state := range serverMonitor {
		if state == server.StatePortClosed {
			slog.Info("Port is closed")
			break
		}
	}
}
