package model

import (
	"testing"
)

func TestDomain(t *testing.T) {

	tree := NewDomainTree()
	tree.AddSubDomain("...com.aaaa.www.............")

	s1 := tree.SubDomain[0]
	if s1.Name != "com" {
		t.Fail()
	}
	s2 := s1.SubDomain[0]
	if s2.Name != "aaaa" {
		t.Fail()
	}
	s3 := s2.SubDomain[0]
	if s3.Name != "www" {
		t.Fail()
	}
	if len(s3.SubDomain) > 0 {
		t.Fail()
	}
}
