package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/logger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	configFile  string
	projectRoot string

	Logger      logger.LoggerConf `yaml:"logger"`
}

func NewConfig(configFile string) Config {
	return Config{
		configFile: configFile,
	}
}

func (c *Config) Parse() error {
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatal(err.Error())
	}

	configYml, err := ioutil.ReadFile(projectRoot + "/" + c.configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configYml, c)
	if err != nil {
		return err
	}

	c.projectRoot = projectRoot
	return nil
}

func getProjectRoot() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		return "", err
	}
	return dir, nil
}
