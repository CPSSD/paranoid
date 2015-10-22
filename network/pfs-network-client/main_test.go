package main

import (
	"testing"
)

func TestParseMessage(t *testing.T) {
	rawWrite := []byte(`{
		"sender": "abc",
		"type": "write",
		"name": "_test.txt",
		"offset": 0,
		"length": 8,
		"data":"aGVsbG8="
	}`)

	rawCreat := []byte(`{
		"sender": "abc",
		"type": "creat",
		"name": "_test.txt"
	}`)

	rawLink := []byte(`{
		"sender": "abc",
		"type": "link",
		"name": "_test.txt",
		"target": "_test2.txt"
	}`)

	rawUnlink := []byte(`{
		"sender": "abc",
		"type": "write",
		"name": "_test.txt"
	}`)

	rawTruncate := []byte(`{
		"sender": "abc",
		"type": "write",
		"name": "_test.txt",
		"offset": 0
	}`)

	// Check is each message parsed without errors
	if _, err := ParseMessage(rawWrite); err != nil {
		t.Error("ParseMessage failed on write")
	}
	if _, err := ParseMessage(rawCreat); err != nil {
		t.Error("ParseMessage failed on creat")
	}
	if _, err := ParseMessage(rawLink); err != nil {
		t.Error("ParseMessage failed on link")
	}
	if _, err := ParseMessage(rawUnlink); err != nil {
		t.Error("ParseMessage failed on unlink")
	}
	if _, err := ParseMessage(rawTruncate); err != nil {
		t.Error("ParseMessage failed on truncate")
	}
}

func TestHasValidFields(t *testing.T) {
	messageWrite := MessageData{"abc123", "write", "_test.txt", "", 0, 8, []byte("aGVsbG8=")}
	messageCreat := MessageData{"abc123", "creat", "_test.txt", "", 0, 0, []byte("")}
	messageLink := MessageData{"abc123", "link", "_test.txt", "_test2.txt", 0, 0, []byte("")}
	messageUnlink := MessageData{"abc123", "unlink", "_test.txt", "", 0, 0, []byte("")}
	messageTruncate := MessageData{"abc123", "truncate", "_test.txt", "", 1, 0, []byte("")}

	// Test if the message has valid fields
	if err := HasValidFields(messageWrite); err != nil {
		t.Error("HasValidFields failed on write")
	}
	if err := HasValidFields(messageCreat); err != nil {
		t.Error("HasValidFields failed on creat")
	}
	if err := HasValidFields(messageLink); err != nil {
		t.Error("HasValidFields failed on link")
	}
	if err := HasValidFields(messageUnlink); err != nil {
		t.Error("HasValidFields failed on unlink")
	}
	if err := HasValidFields(messageTruncate); err != nil {
		t.Error("HasValidFields failed on truncate")
	}
}
