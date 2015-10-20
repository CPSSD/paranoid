package main

import (
	"testing"
)

// Stores the MessageData for each parsed message type
var (
	md1 MessageData
	md2 MessageData
	md3 MessageData
	md4 MessageData
	md5 MessageData
)

func TestParseMessage(t *testing.T) {
	// Write
	fmd1 := []byte(`
		{
			"sender": "abc",
			"type": "write",
			"name": "_test.txt",
			"offset": 0,
			"length": 8,
			"data":"aGVsbG8="
		}
	`)

	// Creat
	fmd2 := []byte(`
		{
			"sender": "abc",
			"type": "creat",
			"name": "_test.txt"
		}
	`)

	// Link
	fmd3 := []byte(`
		{
			"sender": "abc",
			"type": "link",
			"name": "_test.txt",
			"target": "_test2.txt"
		}
	`)

	// Unlink
	fmd4 := []byte(`
		{
			"sender": "abc",
			"type": "write",
			"name": "_test.txt"
		}
	`)

	// truncate
	fmd5 := []byte(`
		{
			"sender": "abc",
			"type": "write",
			"name": "_test.txt",
			"offset": 0
		}
	`)

	// Check is each message parsed without errors
	md1p, err := parseMessage(fmd1)
	if err != nil {
		t.Error("parseMessage failed on write")
	}
	md1 = md1p

	md2p, err := parseMessage(fmd2)
	if err != nil {
		t.Error("parseMessage failed on creat")
	}
	md2 = md2p

	md3p, err := parseMessage(fmd3)
	if err != nil {
		t.Error("parseMessage failed on link")
	}
	md3 = md3p

	md4p, err := parseMessage(fmd4)
	if err != nil {
		t.Error("parseMessage failed on unlink")
	}
	md4 = md4p

	md5p, err := parseMessage(fmd5)
	if err != nil {
		t.Error("parseMessage failed on truncate")
	}
	md5 = md5p
}

func TestHasValidFields(t *testing.T) {
	if err := hasValidFields(md1); err != nil {
		t.Error("hasValidFields failed on write")
	}

	if err := hasValidFields(md2); err != nil {
		t.Error("hasValidFields failed on creat")
	}

	if err := hasValidFields(md3); err != nil {
		t.Error("hasValidFields failed on link")
	}

	if err := hasValidFields(md4); err != nil {
		t.Error("hasValidFields failed on unlink")
	}

	if err := hasValidFields(md5); err != nil {
		t.Error("hasValidFields failed on truncate")
	}
}
