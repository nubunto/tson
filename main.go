package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func listenersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&listeners.listeners); err != nil {
		panic(err)
	}
}

func listenerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Printf("malformed id for listener: %d (%v). Aborting creation.\n", id, err)
		return
	}
	enc := json.NewEncoder(w)
	l := listeners.Find(id)
	if err := enc.Encode(&l); err != nil {
		fmt.Printf("error showing listener for id %d.\n", id)
	}
}

func addListenerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var l listener
	if err := dec.Decode(&l); err != nil {
		panic(err)
	}
	listeners.Start(&l)
}

func removeListenerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Printf("malformed id for listener: %d (%v). Aborting creation.\n", id, err)
		return
	}
	listeners.Remove(id)
}

var listeners *listenerManager

func main() {
	listeners = &listenerManager{
		initPort:  40000,
		listeners: make([]*listener, 0),
	}
	router := httprouter.New()
	router.GET("/listeners", listenersHandler)
	router.GET("/listeners/:id", listenerHandler)
	router.POST("/listeners", addListenerHandler)
	//router.PUT("/listeners/:id", alterListenerHandler)
	router.DELETE("/listeners/:id", removeListenerHandler)
	http.ListenAndServe(":8080", router)
}
