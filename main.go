package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"redis-challenge/internal/server"
	"redis-challenge/internal/store"
)

func main() {
	var err error

	port := flag.Int("port", 6379, "port to listen on")
	useAppendOnlyFile := flag.Bool("aof", false, "use append only file")

	flag.Parse()

	var appendLogReader io.Reader = bytes.NewReader(nil)
	appendLogWriter := io.Discard

	if *useAppendOnlyFile {
		appendLogWriter, err = os.OpenFile("redis-aof.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to open append only file for writing: %v", err))
			os.Exit(1)
		}

		appendLogReader, err = os.Open("redis-aof.log")
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to open append only file for restore: %v", err))
			os.Exit(1)
		}
	}

	serverMonitor := make(server.MonitorChannel)

	srv, err := server.NewChallengeServer(*port, store.NewBuilder()).
		WithWriter(appendLogWriter).
		WithMonitorChannel(serverMonitor).
		RestoreFromArchive(appendLogReader).
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
