package main

import (
	"URLAnalyzer/handler"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Use text format so logs are easy to read in the terminal
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	http.HandleFunc("/", handler.ShowForm)
	http.HandleFunc("/submit", handler.HandleForm)

	slog.Info("Server running", "addr", "http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed to start", "err", err)
	}
}
