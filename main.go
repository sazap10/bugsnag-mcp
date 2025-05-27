package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sazap10/bugsnag-mcp/pkg/config"
	"github.com/sazap10/bugsnag-mcp/pkg/server"
)

const (
	name    = "bugsnag-mcp"
	version = "0.0.1"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Parse command line flags
	transportType := flag.String("transport", "stdio", "Transport type (stdio or sse)")
	sseAddr := flag.String("sse-address", "localhost:8080", "Address for SSE transport")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Set up logging
	level := parseLogLevel(*logLevel)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(consoleHandler)
	slog.SetDefault(logger)

	// Create MCP server
	mcpServer := server.NewMCPServer(name, version, cfg)

	// signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Handle signals
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		slog.Info("Received shutdown signal, shutting down...")
		cancel()
	}()

	// Start the server
	switch *transportType {
	case "stdio":
		slog.Info("Starting bugsnag-mcp with stdio transport")
		if err := server.ServeStdio(ctx, mcpServer); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	case "sse":
		slog.Info("Starting bugsnag-mcp with SSE transport", slog.String("address", *sseAddr))
		if err := server.ServeSSE(ctx, mcpServer, *sseAddr); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	default:
		log.Fatalf("unknown transport type: %s", *transportType)
	}
}

func parseLogLevel(level string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return slog.LevelInfo
	}
	return l
}
