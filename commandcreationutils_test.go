package main

import "testing"

func TestCommandCreationParseUtils(t *testing.T) {

	byteCount := parseByteCount("13")
	if byteCount != 13 {
		t.Errorf("expiryFloat should be 13")
	}
}
