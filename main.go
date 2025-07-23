package main

import (
	"fmt"
	"log/slog"
	"os"
	"redis-challenge/internal/server"
	"redis-challenge/internal/store"
)

func main() {
	dataStore := store.New()
	srv, err := server.NewChallengeServer(0, dataStore)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create server: %v", err))
		os.Exit(1)
	}

	slog.Info("Listening on", "address", srv.Address())

	serverMonitor := make(server.MonitorChannel)
	srv.AddMonitor(serverMonitor)

	for state := range serverMonitor {
		if state == server.StatePortClosed {
			slog.Info("Port is closed")
			break
		}
	}
}
