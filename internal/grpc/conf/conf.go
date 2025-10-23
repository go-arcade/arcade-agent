package conf

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/observabil/arcade/pkg/log"
	"github.com/spf13/viper"
)

type AppConf struct {
	ServerAddr string
	ApiKey     string
	AgentId    string
	cartFile   string
}

var (
	cfg  AppConf
	once sync.Once
)

func NewConf(confDir string) AppConf {
	once.Do(func() {
		var err error
		cfg, err = LoadConfigFile(confDir)
		if err != nil {
			panic(fmt.Sprintf("load conf file error: %s", err))
		}
	})
	return cfg
}

// LoadConfigFile load conf file
func LoadConfigFile(confDir string) (AppConf, error) {

	config := viper.New()
	config.SetConfigFile(confDir) //文件名
	config.AddConfigPath("./conf.d")
	config.SetConfigName("config")
	config.SetConfigType("toml")
	if err := config.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to read configuration file: %v", err)
	}

	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("The configuration changes, re -analyze the configuration file: %s", e.Name)
		if err := config.Unmarshal(&cfg); err != nil {
			_ = fmt.Errorf("failed to unmarshal configuration file: %v", err)
		}
	})
	if err := config.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal configuration file: %v", err)
	}
	fmt.Printf("[Init] config file path: %s\n", confDir)

	return cfg, nil
}
