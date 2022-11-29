package generator

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/webkeydev/websocket"
)

type WsGenerator struct {
	sync.Mutex
	srv      *http.Server
	connChan chan *websocket.Conn
}

func NewWsGenerator() *WsGenerator {
	wsg := &WsGenerator{
		connChan: make(chan *websocket.Conn),
	}
	wsg.listenDummyServer()
	return wsg
}

func (wsg *WsGenerator) listenDummyServer() {
	var wg sync.WaitGroup
	wg.Add(1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", wsg.wsHandler)
	wsg.srv = &http.Server{
		Addr:    ":8999",
		Handler: mux,
	}

	go func() {
		wg.Done()
		_ = wsg.srv.ListenAndServe()
	}()

	wg.Wait()
	return
}

func (wsg *WsGenerator) NewWsConn() (*websocket.Conn, error) {
	wsg.Lock()
	defer wsg.Unlock()

	err := wsClientConnection()
	if err != nil {
		return nil, err
	}
	conn := <-wsg.connChan

	return conn, nil
}

func (wsg *WsGenerator) Release() {
	wsg.Lock()
	defer wsg.Unlock()
	_ = wsg.srv.Close()
}

func (wsg *WsGenerator) wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wsg.connChan <- conn
}

func wsClientConnection() error {
	u := url.URL{Scheme: "ws", Host: "localhost:8999", Path: "/"}
	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return err
}
