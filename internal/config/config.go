package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type DBConfig struct {
	ItemDB string `toml:"item_db"`
	UserDB string `toml:"user_db"`
}

type Config struct {
	DataBase *DBConfig `toml:"database"`
}

func LoadConfigFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	_, err = toml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
