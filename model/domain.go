package model

import (
	"strings"
)

type Domain struct {
	Name      string    `json:"name,omitempty"`
	SubDomain []*Domain `json:"subdomain,omitempty"`
	idx       uint8
}

func NewDomainTree() *Domain {
	tree := new(Domain)
	tree.idx = 0
	tree.Name = "."
	return tree
}

func (d *Domain) AddSubDomain(qdn string) {
	qdn = strings.Trim(qdn, ".")
	idx := strings.Index(qdn, ".")
	if idx < 0 {
		idx = len(qdn)
	}

	subDomain := d.getSubDomain(qdn[:idx])
	if subDomain == nil {
		subDomain = &Domain{
			Name: qdn[:idx],
			idx:  d.idx + 1,
		}

		d.SubDomain = append(d.SubDomain, subDomain)

	}
	if idx < len(qdn) {
		subDomain.AddSubDomain(qdn[idx:])
	}

}

func (d *Domain) getSubDomain(name string) *Domain {
	for _, i := range d.SubDomain {
		if i.Name == name {
			return i
		}
	}
	return nil
}
