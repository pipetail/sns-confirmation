package config

import (
	"flag"
	"fmt"
	"strings"
)

type Application struct {
	// allowed SNS regions
	AllowedRegions  []string
	AllowedAccounts []string
}

func ApplicationFromFlags() (Application, error) {
	cfg := Application{}

	allowedRegions := flag.String("allowed-regions", "", "coma separated list of allowed SNS regions")
	allowedAccounts := flag.String("allowed-accounts", "", "coma separated list of allowed AWS accounts")
	flag.Parse()

	// validate configuration here
	if *allowedRegions == "" {
		return cfg, fmt.Errorf("allowed-regions can't be empty")
	}

	if *allowedAccounts == "" {
		return cfg, fmt.Errorf("allowed-accounts can't be empty")
	}

	// split regions
	cfg.AllowedRegions = strings.Split(*allowedRegions, ",")
	cfg.AllowedAccounts = strings.Split(*allowedAccounts, ",")

	// return the config
	return cfg, nil
}
