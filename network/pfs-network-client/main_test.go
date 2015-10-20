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
	md1p, err := ParseMessage(fmd1)
	if err != nil {
		t.Error("ParseMessage failed on write")
	}
	md1 = md1p

	md2p, err := ParseMessage(fmd2)
	if err != nil {
		t.Error("ParseMessage failed on creat")
	}
	md2 = md2p

	md3p, err := ParseMessage(fmd3)
	if err != nil {
		t.Error("ParseMessage failed on link")
	}
	md3 = md3p

	md4p, err := ParseMessage(fmd4)
	if err != nil {
		t.Error("ParseMessage failed on unlink")
	}
	md4 = md4p

	md5p, err := ParseMessage(fmd5)
	if err != nil {
		t.Error("ParseMessage failed on truncate")
	}
	md5 = md5p
}

func TestHasValidFields(t *testing.T) {
	if err := HasValidFields(md1); err != nil {
		t.Error("HasValidFields failed on write")
	}

	if err := HasValidFields(md2); err != nil {
		t.Error("HasValidFields failed on creat")
	}

	if err := HasValidFields(md3); err != nil {
		t.Error("HasValidFields failed on link")
	}

	if err := HasValidFields(md4); err != nil {
		t.Error("HasValidFields failed on unlink")
	}

	if err := HasValidFields(md5); err != nil {
		t.Error("HasValidFields failed on truncate")
	}
}
