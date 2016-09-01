package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type connInfo struct {
	Addr string `json:"addr"`
	net.Conn
}

type listener struct {
	ID           int `json:"id"`
	Port         int
	net.Listener `json:"-"`
	Out          []*connInfo    `json:"connections"`
	retry        chan *connInfo `json:"-"`
	up           bool
}

var ErrListenerAlreadyUp = errors.New("Listener is already up.")

func (l *listener) start() error {
	if l.up {
		return ErrListenerAlreadyUp
	}
	netListener, err := net.Listen("tcp", fmt.Sprintf(":%d", l.Port))
	if err != nil {
		return err
	}
	l.Listener = netListener
	go l.handleMessages()
	go l.connectOut()
	go l.retryConnections()
	go l.checkConnections()
	l.up = true
	return nil
}

func (l *listener) stop() error {
	log.Printf("tearing listener %d down\n", l.ID)
	return l.Close()
}

func (l *listener) checkConnections() {
	zero := make([]byte, 0)
	for _, c := range l.Out {
		if c.Conn == nil { // shouldn't happen...
			l.retry <- c
			continue
		}
		if _, err := c.Read(zero); err == io.EOF {
			l.retry <- c
		}
	}
}

func (l *listener) retryConnections() {
	for {
		select {
		case out := <-l.retry:
			log.Printf("retrying connection to %s", out.Addr)
			conn, err := net.Dial("tcp", out.Addr)
			if err != nil {
				log.Println("error retrying connection.")
				l.retry <- out
			}
			log.Println("Successful reconnection")
			out.Conn = conn
		}
	}
}

func (l *listener) connectOut() {
	for _, out := range l.Out {
		conn, err := net.Dial("tcp", out.Addr)
		if err != nil {
			l.retry <- out
		}
		out.Conn = conn
	}
}

func (l *listener) handleMessages() {
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		go func() {
			for {
				msg := make([]byte, 256)
				_, err := conn.Read(msg)
				if err != nil {
					break
				}
				fmt.Fprintf(conn, "%s\n", msg)
				fmt.Printf("Got message on %d: %s\n", l.ID, msg)
			}
		}()
	}
}
