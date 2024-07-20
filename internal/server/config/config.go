package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Server address and port
	RunAddr string `env:"ADDRESS" envDefault:"localhost:8080"`
	// Interval between saving metrics to file
	StoreInterval int `env:"STORE_INTERVAL" envDefault:"300"`
	// Path to file with metrics
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	// Restore metrics from file
	Restore bool `env:"RESTORE" envDefault:"true"`
	// Database address
	DatabaseDsn string `env:"DATABASE_DSN" envDefault:""`
	// Key for check sum
	Key string `env:"KEY" envDefault:""`
}

func NewConfig() (Config, error) {
	var cfg Config

	// Setting values by flags, if env not empty, using env
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "interval between saving metrics to file")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path to file with metrics")
	flag.BoolVar(&cfg.Restore, "r", true, "restore metrics from file")
	flag.StringVar(&cfg.DatabaseDsn, "d", "", "database address")
	flag.StringVar(&cfg.Key, "k", "", "key for check sum")

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	flag.Parse()

	// if DatabaseDsn not empty, using it
	addr := strings.Split(cfg.DatabaseDsn, " ")
	if len(addr) > 1 {
		var user string
		var pass string
		var host string
		var port string
		var dbname string

		for _, i := range addr {
			i = strings.Trim(i, `"`)
			variables := strings.Split(i, "=")
			switch variables[0] {
			case "user":
				user = variables[1]
			case "password":
				pass = variables[1]
			case "host":
				host = variables[1]
			case "port":
				port = variables[1]
			case "dbname":
				dbname = variables[1]
			}

		}

		cfg.DatabaseDsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	}

	return cfg, nil
}
