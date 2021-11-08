package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
	internalhttp "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"gopkg.in/yaml.v2"
)

const (
	EnvTest = "test"
	EnvDev  = "dev"  //nolint:deadcode
	EnvProd = "prod" //nolint:deadcode
)

type Config struct {
	configFile  string
	projectRoot string

	Env      string                  `yaml:"env"`
	Server   internalhttp.ServerConf `yaml:"server"`
	Logger   logger.LoggConf         `yaml:"logger"`
	Database interfaces.DatabaseConf `yaml:"database"`
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

func (c *Config) GetProjectRoot() string {
	return c.projectRoot
}

func getProjectRoot() (string, error) {
	_, filename, _, _ := runtime.Caller(0) // nolint
	dir := path.Join(path.Dir(filename), "../../..")
	if err := os.Chdir(dir); err != nil {
		return "", err
	}
	return dir, nil
}
