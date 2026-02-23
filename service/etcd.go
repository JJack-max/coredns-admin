package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/pkg/v3/transport"
)

var client *clientv3.Client

//EtcdInitClient Must first call this function
func EtcdInitClient(conf *config.Config) (err error) {
	etcd := conf.Etcd
	var tlsConfig *tls.Config
	if len(etcd.TLS) == 3 {
		tlsInfo := transport.TLSInfo{
			CertFile:      etcd.TLS[0],
			KeyFile:       etcd.TLS[1],
			TrustedCAFile: etcd.TLS[2],
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			return err
		}
	}
	etcdCfg := clientv3.Config{
		Endpoints: etcd.Endpoint,
		TLS:       tlsConfig,
		Username:  etcd.Username,
		Password:  etcd.Password,
	}

	client, err = clientv3.New(etcdCfg)
	if err != nil {
		return err
	}

	return nil
}

const etcdTimeout = 5 * time.Second

//EtcdGetItems get  etcd items with path as prifix
func EtcdGetItems(path string) (ex []*model.Etcd, err error) {
	kvs, err := EtcdGetKvs(path)
	if err != nil {
		return nil, err
	}

	for k, v := range kvs {

		e := new(model.Etcd)
		e.Key = k

		if err := json.Unmarshal(v, e); err != nil {
			//log.Println(err.Error())
			continue
		}

		ex = append(ex, e)

	}
	return ex, nil
}

//EtcdGetKvs get  etcd keys with path as prifix
func EtcdGetKvs(path string) (kvs map[string][]byte, err error) {
	if client == nil {
		return nil, errors.New("Etcd Client is not initialized")
	}
	kvs = make(map[string][]byte)

	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	r, err := client.Get(ctx, path, clientv3.WithPrefix())
	cancel()

	if err != nil {
		return nil, err
	}

	for _, n := range r.Kvs {
		kvs[string(n.Key)] = n.Value
	}
	return
}

func EtcdPutItems(etcd *model.Etcd) (err error) {

	value, err := json.Marshal(etcd)
	if err != nil {
		return err
	}

	return EtcdPutKv(etcd.Key, string(value))
}

func EtcdPutKv(key, value string) (err error) {
	if client == nil {
		return errors.New("Etcd Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	_, err = client.Put(ctx, key, value)
	cancel()

	return err
}

func EtcdPutKvs(kvs map[string]string, del bool) (err error) {
	if client == nil {
		return errors.New("Etcd Client is not initialized")
	}
	if len(kvs) == 0 {
		return
	}
	var ops []clientv3.Op
	for k, v := range kvs {
		if del {
			idx := strings.LastIndex(k, "/")
			ops = append(ops, clientv3.OpDelete(k[:idx]))
		}

		ops = append(ops, clientv3.OpPut(k, v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	res, err := client.Txn(ctx).If().
		Then(ops...).
		Commit()
	cancel()
	fmt.Println(res)
	return err
}

func EtcdDelete(key string) (err error) {
	if client == nil {
		return errors.New("Etcd Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	_, err = client.Delete(ctx, key)
	cancel()

	return err
}

func EtcdGet(key string) (value []byte, err error) {
	if client == nil {
		return nil, errors.New("Etcd Client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	r, err := client.Get(ctx, key)
	cancel()

	if r.Count != 1 || err != nil {
		return nil, err
	}

	return r.Kvs[0].Value, err
}
