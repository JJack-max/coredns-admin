package service

import (
	"encoding/json"
	"strings"

	"github.com/kexirong/coredns-admin/model"
)

func Domains(path string) (tree *model.Domain, err error) {
	tree = model.NewDomainTree()
	kvs, err := EtcdGetKvs(path)
	if err != nil {
		return nil, err
	}

	for k, v := range kvs {
		if strings.Contains(k, "/arpa/in-addr") || strings.Contains(k, "/dns/ns") {
			continue
		}
		if err := json.Unmarshal(v, new(model.Etcd)); err != nil {
			continue
		}
		key := strings.TrimPrefix(k, path)
		domain := strings.ReplaceAll(key, "/", ".")
		tree.AddSubDomain(domain)

	}
	return tree, nil
}
