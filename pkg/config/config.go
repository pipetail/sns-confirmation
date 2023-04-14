package config

import (
	"flag"
	"fmt"
	"strings"
)

type Application struct {
	// allowed SNS regions
	AllowedRegions []string
	Accounts       []string // TODO: verify accounts
}

func ApplicationFromFlags() (Application, error) {
	cfg := Application{}

	allowedRegions := flag.String("allowed-regions", "", "coma separated list of allowed SNS regions")
	flag.Parse()

	// validate configuration here
	if *allowedRegions == "" {
		return cfg, fmt.Errorf("allowed-regions can't be empty")
	}

	// split regions
	cfg.AllowedRegions = strings.Split(*allowedRegions, ",")

	// return the config
	return cfg, nil
}
