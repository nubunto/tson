package main

import (
	"net"
	"testing"
)

func TestListenerManager(t *testing.T) {
	t.Run("NEXT_PORT", func(t *testing.T) {
		m := &listenerManager{
			initPort: 1,
		}
		if n := m.nextPort(); n != 1 {
			t.Error("initPort is 1, expected 1, got", n)
		}
		if n := m.nextPort(); n != 2 {
			t.Error("initPort is 1 and called once, expected 2, got", n)
		}
	})

	t.Run("ADD", func(t *testing.T) {
		m := &listenerManager{
			initPort: 50000,
		}
		l := &listener{
			Port: m.nextPort(),
		}
		m.Add(l)
	})

	t.Run("REMOVE", func(t *testing.T) {
		m := &listenerManager{
			initPort: 50001,
		}
		l := &listener{
			Port: m.nextPort(),
		}
		m.Add(l)
		if len(m.listeners) != 1 {
			t.Fatalf("Error adding listener to manager: expected 1, got %d", len(m.listeners))
		}
	})

	t.Run("StartAll", func(t *testing.T) {
		m := &listenerManager{
			initPort: 50002,
		}
		l := []*listener{
			&listener{ID: 1},
			&listener{ID: 2},
		}
		m.Add(l...)
		err := m.StartAll()
		if err != nil {
			t.Fatalf("Error starting listeners:", err)
		}

		conn1, err1 := net.Dial("tcp", ":50002")
		if err1 != nil {
			t.Fatalf("Error connecting to listener 50002:", err1)
		}
		conn2, err2 := net.Dial("tcp", ":50003")
		if err2 != nil {
			t.Fatalf("Erro connecting to listener 50003:", err2)
		}
		conn1.Close()
		conn2.Close()
	})

}
