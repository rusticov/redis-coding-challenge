package main

import (
	"fmt"
	"log/slog"
	"os"
	"redis-challenge/internal/server"
	"redis-challenge/internal/store"
)

func main() {
	tracker := store.NewExpiryTracker()
	dataStore := store.New().WithExpiryTracker(tracker)

	scanner := store.NewExpiryScanner(tracker, dataStore)

	serverMonitor := make(server.MonitorChannel)

	srv, err := server.NewChallengeServer(0, dataStore, scanner).
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
