package config

import (
	"github.com/orekhovskiy/shrtn/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	FilePath      string
	DatabaseDSN   string
}

func (c *Config) LogConfig(logger logger.Logger) {
	logger.Info("Application configuration",
		zap.String("server_address", c.ServerAddress),
		zap.String("base_url", c.BaseURL),
		zap.String("file_path", c.FilePath),
		zap.String("database_dsn", c.DatabaseDSN), // sensitive data - log carefully
	)
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
	err = viper.BindEnv("FILE_STORAGE_PATH", "FILE_STORAGE_PATH")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("DATABASE_DSN", "DATABASE_DSN")
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
			if viper.GetString("FILE_STORAGE_PATH") != "" {
				config.FilePath = viper.GetString("FILE_STORAGE_PATH")
			}
			if viper.GetString("DATABASE_DSN") != "" {
				config.DatabaseDSN = viper.GetString("DATABASE_DSN")
			}
		},
	}

	rootCmd.Flags().StringVarP(&config.ServerAddress, "address", "a", "localhost:8080", "HTTP server address (e.g. localhost:8080)")
	rootCmd.Flags().StringVarP(&config.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for shortened URLs (e.g. http://localhost:8080/)")
	rootCmd.Flags().StringVarP(&config.FilePath, "file", "f", "storage.json", "Path to the file for storing URLs (e.g. storage.json)")
	rootCmd.Flags().StringVarP(&config.DatabaseDSN, "database-dsn", "d", "", "Database connection string (e.g. postgres://postgres:password@localhost:5432/shrtn)")
	err = rootCmd.Execute()
	if err != nil {
		return nil, err
	}

	return &config, nil
}
