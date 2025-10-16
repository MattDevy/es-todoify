package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("app", "todoify")

	logger.Info("Hello, world!")
}
