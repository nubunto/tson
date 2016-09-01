package main

import "log"

type listenerManager struct {
	initPort  int
	listeners []*listener
}

func (l *listenerManager) nextPort() int {
	p := l.initPort
	l.initPort++
	return p
}

func (l *listenerManager) Add(listeners ...*listener) {
	l.listeners = append(l.listeners, listeners...)
}

func (l *listenerManager) StartAll() error {
	log.Println("Starting listeners")
	for _, listener := range l.listeners {
		listener.Port = l.nextPort()
		err := listener.start()
		// we don't mind if starting it again
		// not here, at least
		if err != ErrListenerAlreadyUp {
			return err
		}
	}
	log.Println("Done.")
	return nil
}

func (l *listenerManager) Stop(id int) error {
	for _, listener := range l.listeners {
		if listener.ID == id {
			err := listener.stop()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (l *listenerManager) Remove(id int) error {
	log.Println("Removing listener", id)
	for i, listener := range l.listeners {
		if listener.ID == id {
			err := listener.stop()
			if err != nil {
				return err
			}
			l.listeners[i] = nil // garbage collection
			l.listeners = append(l.listeners[:i], l.listeners[i+1:]...)
		}
	}
	log.Println("Done.")
	return nil
}

func (l *listenerManager) Find(id int) *listener {
	for _, listener := range l.listeners {
		if listener.ID == id {
			return listener
		}
	}
	return nil
}
