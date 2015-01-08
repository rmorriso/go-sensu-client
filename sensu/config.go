package sensu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"path/filepath"
)

type ClientConfig struct {
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	Subscriptions []string `json:"subscriptions"`
}

type RabbitmqConfigSSL struct {
	PrivateKeyFile string `json:"private_key_file"`
	CertChainFile  string `json:"cert_chain_file"`
}

type RabbitmqConfig struct {
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Vhost    string            `json:"vhost"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	Ssl      RabbitmqConfigSSL `json:"ssl"`
}

type Config struct {
	Checks   map[string]Check `json:"checks"`
	Client   ClientConfig     `json:"client"`
	Rabbitmq RabbitmqConfig   `json:"rabbitmq"`
	rawData  *simplejson.Json
}

func LoadConfigs(configFile string, configDirs []string) (*Config, error) {
	js, ferr := parseFile(configFile)
	if ferr != nil {
		log.Printf("Unable to open config file: %s", ferr)
	}

	for _, dir := range configDirs {
		files, derr := ioutil.ReadDir(dir)
		if derr != nil {
			log.Printf("Unable to open config directory: %s", derr)
		}

		for _, f := range files {
			jsd, err := parseFile(filepath.Join(dir, f.Name()))
			if err != nil {
				log.Printf("Could not load %s: %s", f.Name(), err)
				continue
			}

			err = js.Extend(jsd)
			if err != nil {
				log.Printf("Error merging configs: %s", err)
			}
		}
	}

	//Reencoding merged JSON to parse to concrete type
	if nil == js {
		return nil, errors.New("There was no configuration.")
	}
	mergedJson, err := json.Marshal(js.data)
	if err != nil {
		return nil, errors.New("Unable to reencode merged json")
	}
	config := new(Config)
	json.Unmarshal(mergedJson, &config)
	config.rawData, _ = simplejson.NewJson(mergedJson)

	validationErrors := validateConfig(config)
	if len(validationErrors) > 0 {
		for _, e := range validationErrors {
			fmt.Print(e)
		}
		return nil, errors.New("Error validating the configs")
	}

	return config, nil
}

func parseFile(filename string) (*Json, error) {
	j := new(Json)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("File error: %v", err)
	}

	err = json.Unmarshal(file, &j.data)
	if err != nil {
		return nil, fmt.Errorf("json error: %v", err)
	}

	return j, nil
}

func validateConfig(cfg *Config) []error {
	errs := []error{}

	if cfg.Client.Address == "" {
		errs = append(errs, errors.New("Missing client address"))
	}

	return errs
}

func (c *Config) Data() *simplejson.Json {
	return c.rawData
}
