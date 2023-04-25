package config

import (
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Logger *logger.Config `mapstructure:"logger"`
	DB     *db.Config     `mapstructure:"db"`
	OOB    *OOB           `mapstructure:"oob"`
	Web    *Web           `mapstructure:"web"`
}

type OOB struct {
	DNS  DNS  `mapstructure:"dns"`
	JNDI JNDI `mapstructure:"jndi"`
}

type DNS struct {
	Main string   `mapstructure:"main"`
	NS   []string `mapstructure:"ns"`
	IP   string   `mapstructure:"ip"`
}

type JNDI struct {
	Address string `mapstructure:"address"`
	Limit   int64  `mapstructure:"limit"`
}

type Web struct {
	Address string `mapstructure:"address"`
	Github  struct {
		ClientID     string `mapstructure:"client-id"`
		ClientSecret string `mapstructure:"client-secret"`
	} `mapstructure:"github"`
}

func Load(path string) (v *Config, err error) {
	viper.SetConfigFile(path)

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&v); err != nil {
		return nil, err
	}

	return
}
