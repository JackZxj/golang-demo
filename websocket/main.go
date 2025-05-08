package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var status bool
var cache = make(chan []byte, 1024)
var consumers = map[string]func(msg []byte){}
var rwLock sync.RWMutex
var onUpload uint32

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func addConsumer(id string, consumer func(msg []byte)) {
	rwLock.Lock()
	defer rwLock.Unlock()
	consumers[id] = consumer
}

func removeConsumer(id string) {
	rwLock.Lock()
	defer rwLock.Unlock()
	delete(consumers, id)
}

func listConsumers() []func(msg []byte) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	result := make([]func(msg []byte), 0, len(consumers))
	for _, consumer := range consumers {
		result = append(result, consumer)
	}
	return result
}

func pong(msgType int, msg []byte, conn *websocket.Conn, id string) bool {
	if msgType == websocket.TextMessage && strings.Contains(string(msg), "ping") { // 心跳检测，{"type":"ping"}
		err := conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		if err != nil {
			log.Printf("[%s] response pong error %v", id, err)
		}
		return true
	}
	return false
}

func readStreamHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	id := fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().Unix())
	addConsumer(id, func(msg []byte) {
		if conn == nil || conn.NetConn() == nil {
			removeConsumer(id)
			return
		}
		err = conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			log.Printf("[%s] Error writing message: %v", id, err)
		}
	})
	defer removeConsumer(id)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received: %s\n", msg)
		if pong(msgType, msg, conn, id) {
			continue
		}
	}
}

func uploadStreamHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadUint32(&onUpload) == 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already on uploading"))
		return
	}
	atomic.StoreUint32(&onUpload, 1)
	defer atomic.StoreUint32(&onUpload, 0)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set websocket upgrade"))
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received: %d\n", len(msg))
		if pong(msgType, msg, conn, "uploader") {
			continue
		}

		if msgType == websocket.BinaryMessage {
			cache <- msg
		}
	}

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if status {
			w.Write([]byte("on"))
		} else {
			w.Write([]byte("off"))
		}
	case http.MethodPost:
		result := r.URL.Query().Get("status")
		status = result == "on"
		w.Write([]byte("ok"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func consume() {
	for {
		msg, ok := <-cache
		if !ok {
			panic("channel closed")
		}

		var wg sync.WaitGroup
		for _, consumer := range listConsumers() {
			wg.Add(1)
			go func(consumer func([]byte)) {
				defer wg.Done()
				if e := recover(); e != nil {
					log.Printf("consumer error: %v", e)
				}
				consumer(msg)
			}(consumer)
		}
		wg.Wait()
	}
}

func saveFrame(data []byte) {
	filename := fmt.Sprintf("images/frame_%d.jpg", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Printf("Failed to write to file: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/upload-stream", uploadStreamHandler)
	http.HandleFunc("/read-stream", readStreamHandler)
	http.HandleFunc("/", webHandler)

	status = true // default on

	log.Println("Starting server at port 8080")
	// addConsumer("save", saveFrame)
	go consume()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
