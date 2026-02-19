package main

import (
	"fmt"
	"log"

	"github.com/devaloi/configo"
)

type DatabaseConfig struct {
	Host     string `config:"database.host" default:"localhost" validate:"required"`
	Port     int    `config:"database.port" default:"5432" validate:"required,min=1,max=65535"`
	Name     string `config:"database.name" validate:"required"`
	Password string `config:"database.password" default:"secret"`
}

func main() {
	cfg := configo.New(
		configo.WithFile("config.yaml"),
		configo.WithEnvPrefix("APP"),
	)

	if err := cfg.Load(); err != nil {
		log.Fatal(err)
	}

	var db DatabaseConfig
	if err := cfg.Bind(&db); err != nil {
		log.Fatal(err)
	}

	if err := cfg.ValidateStruct(db); err != nil {
		log.Printf("Validation warnings: %v", err)
	}

	fmt.Printf("Database: %s:%d/%s\n", db.Host, db.Port, db.Name)
}
