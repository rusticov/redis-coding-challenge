package config

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Configuration struct {
	Port            int
	AppendLogReader io.Reader
	AppendLogWriter io.Writer
}

func LoadConfiguration() Configuration {
	configuration := Configuration{
		Port:            6379,
		AppendLogReader: bytes.NewReader(nil),
		AppendLogWriter: io.Discard,
	}

	flag.IntVar(&configuration.Port, "port", 6379, "port to listen on")
	useAppendOnlyFile := flag.Bool("aof", false, "use append only file")

	flag.Parse()

	if *useAppendOnlyFile {
		var err error

		configuration.AppendLogWriter, err = os.OpenFile("redis-aof.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to open append only file for writing: %v", err))
			os.Exit(1)
		}

		configuration.AppendLogReader, err = os.Open("redis-aof.log")
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to open append only file for restore: %v", err))
			os.Exit(1)
		}
	}
	return configuration
}
