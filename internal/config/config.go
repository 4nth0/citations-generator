package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFilPath = "generator.yaml"
)

type Config struct {
	Base   string `yaml:"base"`
	Author struct {
		Name string `yaml:"name"`
		Bio  string `yaml:"bio"`
	} `yaml:"author"`
	Meta struct {
		OG struct {
			Title string `yaml:"title"`
		} `yaml:"og"`
	} `yaml:"meta"`
	Generator struct {
		CitationsPerPage int    `yaml:"citations_per_page"`
		Source           string `yaml:"source"`
		Paths            struct {
			Index   string `yaml:"index"`
			Detail  string `yaml:"detail"`
			Listing string `yaml:"listing"`
		} `yaml:"paths"`
		Templates struct {
			Index    Template          `yaml:"index"`
			Detail   Template          `yaml:"detail"`
			Listing  Template          `yaml:"listing"`
			Partials map[string]string `yaml:"partials"`
		} `yaml:"templates"`
	} `yaml:"generator"`
}

type Template struct {
	Template string `yaml:"template"`
	Dest     string `yaml:"dest"`
}

func Load() (*Config, error) {

	file, err := os.Open(ConfigFilPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
