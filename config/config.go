package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

//var filePath = "config.yaml"

//Config all config in this
type Config struct {
	Host         string `yaml:"host,omitempty"`
	Port         string `yaml:"port,omitempty"`
	UserEtcdPath string `yaml:"user_etcd_path,omitempty"`
	Etcd         struct {
		Endpoint   []string `yaml:"endpoint"`
		PathPrefix string   `yaml:"path_prefix"`
		Username   string   `yaml:"username"`
		Password   string   `yaml:"password"`
		Timeout    int      `yaml:"timeout"`
		TLS        []string `yaml:"tls"`
	}
}

var conf *Config

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func loadConfig(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}

//Get  get a congfig instance
func Get() *Config {
	if conf == nil {
		log.Fatalln("Configuration is not initialized")
	}
	return conf
}

func Set(filePath string) {

	if !pathIsExist(filePath) {
		log.Fatalln("The configuration file does not exist")
	}

	conf = new(Config)

	loadConfig(filePath)
	loadDefaultConfig(conf)

}

func LoadDefaultConfig() {
	conf = new(Config)
	loadDefaultConfig(conf)
}

// func init() {
// 	Get()
// }
