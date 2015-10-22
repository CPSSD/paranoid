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
	if _, err := parseMessage(rawWrite); err != nil {
		t.Error("parseMessage failed on write")
	}
	if _, err := parseMessage(rawCreat); err != nil {
		t.Error("parseMessage failed on creat")
	}
	if _, err := parseMessage(rawLink); err != nil {
		t.Error("parseMessage failed on link")
	}
	if _, err := parseMessage(rawUnlink); err != nil {
		t.Error("parseMessage failed on unlink")
	}
	if _, err := parseMessage(rawTruncate); err != nil {
		t.Error("parseMessage failed on truncate")
	}
}

func TestHasValidFields(t *testing.T) {
	x, y := 0, 8
	messageWrite := MessageData{"abc123", "write", "_test.txt", "", &x, &y, []byte("aGVsbG8=")}
	messageCreat := MessageData{"abc123", "creat", "_test.txt", "", nil, nil, []byte("")}
	messageLink := MessageData{"abc123", "link", "_test.txt", "_test2.txt", nil, nil, []byte("")}
	messageUnlink := MessageData{"abc123", "unlink", "_test.txt", "", nil, nil, []byte("")}
	messageTruncate := MessageData{"abc123", "truncate", "_test.txt", "", nil, nil, []byte("")}

	// Test if the message has valid fields
	if err := hasValidFields(messageWrite); err != nil {
		t.Error("hasValidFields failed on write")
	}
	if err := hasValidFields(messageCreat); err != nil {
		t.Error("hasValidFields failed on creat")
	}
	if err := hasValidFields(messageLink); err != nil {
		t.Error("hasValidFields failed on link")
	}
	if err := hasValidFields(messageUnlink); err != nil {
		t.Error("hasValidFields failed on unlink")
	}
	if err := hasValidFields(messageTruncate); err != nil {
		t.Error("hasValidFields failed on truncate")
	}
}
