package utils

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

type Configuration struct {
	Debug    bool   `json:"debug"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Url      string `json:"url"`
	Database struct {
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"database"`
	Assets struct {
		Static    string `json:"static"`
		Templates string `json:"templates"`
	} `json:"assets"`
}

var config *Configuration

func Config() *Configuration {
	if config == nil {
		config = &Configuration{}

		viper.SetConfigName("config")
		viper.AddConfigPath(fmt.Sprintf("$%s/", viper.Get("CONFIG")))
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			flag.Parse() // flags are not yet parsed, avoid error from glog
			glog.Fatalf("when reading config file: %s", err)
		}
		err = viper.Unmarshal(config)
		if err != nil {
			flag.Parse() // flags are not yet parsed, avoid error from glog
			glog.Fatalf("when unmarshalling the json: %s", err)
		}
	}
	return config
}

func init() {
	viper.Set("PROJECT", "github.com/alienantfarm/anthive")
	viper.Set("CONFIG", "ANTHIVE_CONFIG")
}
