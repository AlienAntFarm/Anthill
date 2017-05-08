package utils

import (
	"encoding/json"
	"github.com/golang/glog"
	"os"
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

var Config = getConfig()

func getConfig() (config Configuration) {
	configFile, err := os.Open(getConfigPath())
	if err != nil {
		glog.Fatalf("when reading config file: %s", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		glog.Fatalf("when unmarshalling the json: %s", err)
	}
	return
}

func getConfigPath() (configPath string) {
	if configPath = os.Getenv(CONFIG); configPath == "" {
		configPath = "config.json"
	}
	return
}
