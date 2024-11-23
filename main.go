package main

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ashihara-api/core/utils/logger"

	"github.com/ashihara-api/geo/interface/transport/http/server"
)

func convertLogLevel(lv string) (slog.Level, error) {
	switch strings.ToLower(lv) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	}
	return 0, fmt.Errorf("%s is unsupported log level", lv)
}

func main() {
	var (
		listenAddr string
		logLvStr   string
	)
	flag.StringVar(&listenAddr, "listen-addr", ":80", "server listen address")
	flag.StringVar(&logLvStr, "log-level", "info", "log level")
	flag.Parse()

	lv, err := convertLogLevel(logLvStr)
	if err != nil {
		panic(err)
	}

	conf := server.Config{
		Address: listenAddr,
		Logger:  logger.NewWithLevel(lv),
	}
	server.Run(conf)
}
