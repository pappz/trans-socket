package wsgenerator

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/webkeydev/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mutex = &sync.Mutex{}

	srv *http.Server

	connChan = make(chan *websocket.Conn)
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connChan <- conn
}

func wsClientConnection() {
	u := url.URL{Scheme: "ws", Host: "localhost:8999", Path: "/"}
	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("dial: %s", err)
	}

}

func initDummyServer() {
	var wg sync.WaitGroup
	wg.Add(1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", wsHandler)
	srv = &http.Server{
		Addr:    ":8999",
		Handler: mux,
	}

	go func() {
		wg.Done()
		_ = srv.ListenAndServe()
	}()

	wg.Wait()
	return
}

func NewWsConn() *websocket.Conn {
	mutex.Lock()
	defer mutex.Unlock()

	if srv == nil {
		initDummyServer()
	}

	wsClientConnection()
	conn := <-connChan

	return conn
}

func Release() {
	mutex.Lock()
	defer mutex.Unlock()
	_ = srv.Close()
	srv = nil
}
