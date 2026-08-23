package config

import "github.com/spf13/viper"

type Config struct {
	SidecarPort int
	DaemonPort  int
	LogLevel     string
}

func Load(cfgFile string) *Config {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".go-daemon-template")
	}

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DAEMON_PORT", "2222")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	// Config file is optional
	_ = viper.ReadInConfig()

	return &Config{
		SidecarPort: viper.GetInt("PORT"),
		DaemonPort:  viper.GetInt("DAEMON_PORT"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
	}
}