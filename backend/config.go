package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`

	API     APIConfig     `yaml:"api"`
	Usecase UsecaseConfig `yaml:"usecase"`
	DB      DBConfig      `yaml:"db"`
	JWT     JWTConfig     `yaml:"jwt"`
}

type APIConfig struct {
	MinPasswordSize int `yaml:"min_password_size"`
	MaxPasswordSize int `yaml:"max_password_size"`
	MinUsernameSize int `yaml:"min_username_size"`
	MaxUsernameSize int `yaml:"max_username_size"`
}

type UsecaseConfig struct {
	DefaultMessage string `yaml:"default_message"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

func LoadConfig(pathToFile string) (*Config, error) {
	filename, err := filepath.Abs(pathToFile)
	if err != nil {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
