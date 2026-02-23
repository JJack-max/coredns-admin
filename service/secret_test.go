package service

import "testing"

func TestSecret(t *testing.T) {
	password := "password"
	secret := MakeSecret(password)
	if VerifySecret(secret, password) {
		t.Fatal(secret, password)
	}
}
