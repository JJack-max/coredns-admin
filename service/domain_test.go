package service

import (
	"testing"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
)

func TestDomains(t *testing.T) {
	config.LoadDefaultConfig()
	err := EtcdInitClient(config.Get())
	if err != nil {
		t.Error(err)
	}
	tree, err := Domains("/coredns/")
	if err != nil {
		t.Error(err)
	}
	var f func(d *model.Domain)
	f = func(d *model.Domain) {

		for _, v := range d.SubDomain {
			if len(v.SubDomain) > 0 {
				f(v)
			}
		}
	}
	f(tree)

}
