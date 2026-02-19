package main

import (
	"fmt"
	"log"

	"github.com/devaloi/configo"
)

func main() {
	cfg := configo.New(
		configo.WithDefaults(map[string]any{
			"server": map[string]any{
				"host": "localhost",
				"port": 3000,
			},
		}),
		configo.WithFile("config.yaml"),
		configo.WithEnvPrefix("APP"),
	)

	if err := cfg.Load(); err != nil {
		log.Fatal(err)
	}

	host := configo.GetOr[string](cfg, "server.host", "localhost")
	port := configo.GetOr[int](cfg, "server.port", 3000)
	debug := configo.GetOr[bool](cfg, "server.debug", false)

	fmt.Printf("Server: %s:%d (debug=%v)\n", host, port, debug)
}
