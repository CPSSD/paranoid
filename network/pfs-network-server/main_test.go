// +build !integration
package main

import (
	"testing"
	"time"
)

func TestSendReceiveMessage(t *testing.T) {
	pool := NewConnPool()
	go pool.Run()
	conn := NewConnection(nil, pool)
	pool.Register <- conn
	defer func() { pool.Unregister <- conn }()
	pool.messages <- []byte("test!")
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	select {
	case message := <-conn.messages:
		if string(message) != "test!" {
			t.Error("Output text didn't match input. Expected: test!. Actual: ",
				string(message))
		}
	case <-timeout:
		t.Error("Message not received.")
	}
}
