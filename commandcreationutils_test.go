package main

import "testing"

func TestCommandCreationParseUtils(t *testing.T) {

	expiryFloat := parseExpiry("1")
	if expiryFloat != 1 {
		t.Errorf("expiryFloat should be 1")
	}
	byteCount := parseByteCount("13")
	if byteCount != 13 {
		t.Errorf("expiryFloat should be 13")
	}
}
