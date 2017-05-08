package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

var (
	Config  = getConfig()
	Trace   = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info    = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error   = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func getConfig() (config Configuration) {
	configFile, err := os.Open(getConfigPath())
	if err != nil {
		Error.Fatalf("when reading config file: %s", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		Error.Fatalf("when unmarshalling the json: %s", err)
	}
	if !config.Debug {
		Trace = log.New(ioutil.Discard, "", 0)
	}
	return
}

func getConfigPath() (configPath string) {
	if configPath = os.Getenv(CONFIG); configPath == "" {
		configPath = "config.json"
	}
	return
}
