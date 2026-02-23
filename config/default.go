package config

import "strings"

const (
	defaultEndpoint     = "http://10.1.1.224:2379"
	defaultPathPrefix   = "/coredns"
	defaultUserEtcdPath = "/user/coredns"
	defaultEtcdTimeout  = 5
	defaultHost         = ""
	defaultPort         = "8080"
)

func loadDefaultConfig(conf *Config) {

	if conf.Host == "" {
		conf.Host = defaultHost
	}
	if conf.Port == "" {
		conf.Port = defaultPort
	}
	if conf.UserEtcdPath == "" {
		conf.UserEtcdPath = defaultUserEtcdPath
	}
	if !strings.HasSuffix(conf.UserEtcdPath, "/") {
		conf.UserEtcdPath += "/"
	}
	if len(conf.Etcd.Endpoint) == 0 {
		conf.Etcd.Endpoint = []string{defaultEndpoint}
	}
	if conf.Etcd.PathPrefix == "" {
		conf.Etcd.PathPrefix = defaultPathPrefix
	}
	if !strings.HasSuffix(conf.Etcd.PathPrefix, "/") {
		conf.Etcd.PathPrefix += "/"
	}
	if conf.Etcd.Timeout == 0 {
		conf.Etcd.Timeout = defaultEtcdTimeout
	}

}
