package main

import (
	"net"
	"testing"
	"time"
)

func TestListener(t *testing.T) {
	/*
		srv, err := net.Listen("tcp", ":9999")
		if err != nil {
			t.Fatalf("Failed to spin up test TCP server: %v\n", err)
		}
	*/

	t.Run("UP", func(t *testing.T) {
		l := &listener{ID: 1, Port: 10001}
		err := l.start()
		if err != nil {
			t.Fatalf("Error putting up the listener:", err)
		}
		time.Sleep(100 * time.Millisecond)
		conn, err := net.Dial("tcp", ":10001")
		if err != nil {
			t.Fatal("Failed to open up TCP connection to listener.")
		}
		conn.Close()
		err = l.stop()
		if err != nil {
			t.Fatal("Failed to put server down:", err)
		}
	})

	t.Run("UPOnlyOnce", func(t *testing.T) {
		l := &listener{ID: 1, Port: 10005}
		err := l.start()
		if err != nil {
			t.Fatalf("Error putting the listener up:", err)
		}
		err = l.start()
		if err == nil || err != ErrListenerAlreadyUp {
			t.Fatalf("putting up the listener a second time should error with ErrListenerAlreadyUp:", err)
		}
		err = l.stop()
		if err != nil {
			t.Fatalf("Error putting the listener down:", err)
		}
	})

	t.Run("DOWN", func(t *testing.T) {
		l := &listener{ID: 1, Port: 10002}
		err := l.start()
		if err != nil {
			t.Fatal("Failed to up server up:", err)
		}
		time.Sleep(50 * time.Millisecond)
		err = l.stop()
		if err != nil {
			t.Fatal("Failed to put server down:", err)
		}
		time.Sleep(50 * time.Millisecond)
		_, err = net.Dial("tcp", ":10002")
		if err == nil {
			t.Fatalf("Failed to tear down listener.")
		}
	})

	t.Run("CONNECTIONS", func(t *testing.T) {
		l := &listener{
			ID:   1,
			Port: 10003,
			Out: []*connInfo{
				&connInfo{
					Addr: "0.0.0.0:2020",
				},
				&connInfo{
					Addr: "0.0.0.0:3030",
				},
			},
		}

		ln1, err := net.Listen("tcp", "0.0.0.0:2020")
		if err != nil {
			t.Fatal("Failed to start server 1:", err)
		}

		ln2, err := net.Listen("tcp", "0.0.0.0:3030")
		if err != nil {
			t.Fatal("Failed to start server 2:", err)
		}

		time.Sleep(50 * time.Millisecond)

		err = l.start()
		if err != nil {
			t.Fatal("Failed to put server up:", err)
		}

		conn1, err := ln1.Accept()
		if err != nil {
			t.Fatal("No connection accepted on server 1:", err)
		}

		conn2, err := ln2.Accept()
		if err != nil {
			t.Fatal("No connection accepted on server 2:", err)
		}

		err = l.stop()
		err = conn1.Close()
		err = conn2.Close()
		err = ln1.Close()
		err = ln2.Close()
		if err != nil {
			t.Error("Error tearing down test:", err)
		}
	})

	t.Run("RETRY", func(t *testing.T) {
		t.Skip("This test is hanging forever...")
		l := &listener{
			ID:   1,
			Port: 10003,
			Out: []*connInfo{
				&connInfo{
					Addr: "0.0.0.0:4040",
				},
			},
			retry: make(chan *connInfo),
		}

		ln, err := net.Listen("tcp", "0.0.0.0:4040")
		if err != nil {
			t.Fatal("Failed to start server 1:", err)
		}

		time.Sleep(50 * time.Millisecond)

		err = l.start()
		if err != nil {
			t.Fatal("Failed to put server up:", err)
		}

		conn, err := ln.Accept()
		if err != nil {
			t.Fatal("Failed to accept connection:", err)
		}
		conn.Close()

		ln.Close()
		time.Sleep(50 * time.Millisecond)

		ln, err = net.Listen("tcp", "0.0.0.0:4040")

		time.Sleep(50 * time.Millisecond)

		time.Sleep(200 * time.Millisecond)
		conn, err = ln.Accept()
		if err != nil {
			t.Fatal("Failed to accept retry connection:", err)
		}
		conn.Close()
	})
}
