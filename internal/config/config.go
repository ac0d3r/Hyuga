package config

import (
	"hyuga/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	DebugMode bool          `mapstructure:"debug"`
	Logger    logger.Config `mapstructure:"logger"`
	Db        DB            `mapstructure:"db"`
	OOB       OOB           `mapstructure:"oob"`
	Api       Api           `mapstructure:"api"`
}

type DB struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

type OOB struct {
	Dns struct {
		Domain string   `mapstructure:"domain"`
		NS     []string `mapstructure:"ns"`
		IP     string   `mapstructure:"ip"`
	} `mapstructure:"dns"`
	Jndi struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"jndi"`
}

type Api struct {
	Address string `mapstructure:"address"`
}

func Load(path string) (v *Config, err error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&v); err != nil {
		return nil, err
	}

	if v.DebugMode {
		v.Logger.Level = "debug"
	} else {
		v.Logger.Level = "info"
	}

	return
}
