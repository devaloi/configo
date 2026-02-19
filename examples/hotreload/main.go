package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devaloi/configo"
)

func main() {
	cfg := configo.New(
		configo.WithFile("config.yaml"),
	)

	if err := cfg.Load(); err != nil {
		log.Fatal(err)
	}

	cfg.OnChange(func(c *configo.Config) {
		host := configo.GetOr[string](c, "server.host", "unknown")
		port := configo.GetOr[int](c, "server.port", 0)
		fmt.Printf("Config reloaded: %s:%d\n", host, port)
	})

	if err := cfg.Watch(); err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cfg.StopWatch() }()

	fmt.Println("Watching config.yaml for changes. Press Ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down.")
}
