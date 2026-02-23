package config

import (
	"testing"
)

func TestGet(t *testing.T) {

}

func TestLoadDefaultConfig(t *testing.T) {
	LoadDefaultConfig()
	conf := Get()
	if conf.Host != defaultHost {
		t.Fatalf("conf.Host!=defaultHost")
	}
	if conf.Port != defaultPort {
		t.Fatalf("conf.Port != defaultPort")
	}
	etcd := conf.Etcd
	// loadDefaultConfig appends trailing slash to PathPrefix
	if etcd.PathPrefix != defaultPathPrefix+"/" {
		t.Fatalf("etcd.PathPrefix != defaultPathPrefix+/")
	}
	if etcd.Endpoint[0] != defaultEndpoint {
		t.Fatalf("etcd.Endpoint[0] != defaultEndpoint")
	}
	if etcd.Timeout != defaultEtcdTimeout {
		t.Fatalf("etcd.Timeout != defaultEtcdTimeout ")
	}
}

func TestLoadConfig(t *testing.T) {
	loadConfig("../config.yaml")

}
