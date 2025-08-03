package config

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

type Configuration struct {
	Port            int
	AppendLogReader io.Reader
	AppendLogWriter io.Writer
}

func LoadConfiguration() (Configuration, error) {
	configuration := Configuration{
		Port:            6379,
		AppendLogReader: bytes.NewReader(nil),
		AppendLogWriter: io.Discard,
	}
	useAppendOnlyFile := false

	flag.IntVar(&configuration.Port, "port", 6379, "port to listen on")
	flag.BoolVar(&useAppendOnlyFile, "aof", false, "use append only file")

	flag.Parse()

	if useAppendOnlyFile {
		var err error

		configuration.AppendLogWriter, err = os.OpenFile("redis-aof.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return Configuration{}, fmt.Errorf("failed to open append only file for writing: %w", err)
		}

		configuration.AppendLogReader, err = os.Open("redis-aof.log")
		if err != nil {
			return Configuration{}, fmt.Errorf("failed to open append only file for restore: %w", err)
		}
	}

	return configuration, nil
}
