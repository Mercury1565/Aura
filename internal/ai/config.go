package ai

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mercury1565/Aura/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	GroqAPIKey      string `mapstructure:"groq_api_key"`
	GeminiAPIKey    string `mapstructure:"gemini_api_key"`
	ModelName       string `mapstructure:"model_name"`
	BaseInstruction string `mapstructure:"base_instruction"`
}

func LoadConfig() (*Config, error) {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "aura")
	configPath := filepath.Join(configDir, "config.yaml")

	// Explicitly set the config file path
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Try to read the file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !os.IsNotExist(err) {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set default BaseInstruction if it's empty
	if cfg.BaseInstruction == "" {
		cfg.SetAndSave("base_instruction", utils.BaseInstruction)
	}

	return &cfg, nil
}

func (c *Config) SetAndSave(key string, value string) error {
	viper.Set(key, value)

	// Ensure the directory exists before saving
	configDir := filepath.Dir(viper.ConfigFileUsed())
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	// WriteConfig writes the current state to the path set in viper.SetConfigFile
	// or the path discovered during LoadConfig
	return viper.WriteConfigAs(viper.ConfigFileUsed())
}

func (c *Config) HandleConfigSet(osArgs []string) {
	key := os.Args[2]
	value := os.Args[3]

	err := c.SetAndSave(key, value)
	if err != nil {
		fmt.Printf("Failed to update config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s to: %s\n", key, value)
	fmt.Printf("📂 Config saved at: %s\n", viper.ConfigFileUsed())
}

func (c *Config) Validate() error {
	if c.GroqAPIKey == "" {
		return fmt.Errorf("Error validating config file at %s: missing groq_api_key", viper.ConfigFileUsed())
	}
	if c.ModelName == "" {
		return fmt.Errorf("Error validating config file at %s: missing model_name", viper.ConfigFileUsed())
	}
	return nil
}
