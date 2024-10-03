package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func InitializeConfig() (*Config, error) {
	var config Config

	viper.AutomaticEnv()
	err := viper.BindEnv("SERVER_ADDRESS", "SERVER_ADDRESS")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("BASE_URL", "BASE_URL")
	if err != nil {
		return nil, err
	}
	var rootCmd = &cobra.Command{
		Use:   "shrtn",
		Short: "URL Shortener Service",
		Run: func(cmd *cobra.Command, args []string) {
			// Get config is such priority ENV > flags > default values
			if viper.GetString("SERVER_ADDRESS") != "" {
				config.ServerAddress = viper.GetString("SERVER_ADDRESS")
			}
			if viper.GetString("BASE_URL") != "" {
				config.BaseURL = viper.GetString("BASE_URL")
			}
		},
	}

	rootCmd.Flags().StringVarP(&config.ServerAddress, "address", "a", "localhost:8080", "HTTP server address (e.g. localhost:8080)")
	rootCmd.Flags().StringVarP(&config.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for shortened URLs (e.g. http://localhost:8080/)")

	err = rootCmd.Execute()
	if err != nil {
		return nil, err
	}

	return &config, nil
}
