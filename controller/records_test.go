package controller

import "testing"

func TestGrowBasicPrefix(t *testing.T) {
	if growBasicPrefix("") != "#1" {
		t.Fatal()
	}
	if growBasicPrefix("#") != "#1" {
		t.Fatal()
	}
	if growBasicPrefix("#1") != "#2" {
		t.Fatal()
	}
	if growBasicPrefix("#19") != "#191" {
		t.Fatal()
	}
}
