package config

import (
	"github.com/spf13/cobra"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func InitializeConfig() (*Config, error) {
	var config Config

	var rootCmd = &cobra.Command{
		Use:   "shrtn",
		Short: "URL Shortener Service",
	}

	rootCmd.Flags().StringVarP(&config.ServerAddress, "address", "a", "localhost:8080", "HTTP server address (e.g. localhost:8080)")
	rootCmd.Flags().StringVarP(&config.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for shortened URLs (e.g. http://localhost:8080/)")

	err := rootCmd.Execute()
	if err != nil {
		return nil, err
	}

	return &config, nil
}
