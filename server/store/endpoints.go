package store

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
	"fmt"
)

var (
	n = 0
	mu = &sync.Mutex{}
	clients = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{ReadBufferSize: 1024,WriteBufferSize: 1024}
)


// create new endpoint factory
func NewEndpointsFactory(messageStore MessageStore) *endpointsFactory {
	return &endpointsFactory{messageStore: messageStore}
}

// endpointsFactory
type endpointsFactory struct {
	messageStore MessageStore
}



// endpoint for websocket
func(ef *endpointsFactory) WebsocketEndpoint() func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		n++
		upgrader.CheckOrigin = func(r *http.Request) bool {return true}
		conn, err := upgrader.Upgrade(w,r,nil)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		mu.Lock()
		clients[conn]=true
		mu.Unlock()
		//get all messages from mongo and and send to client
		msgs,err := ef.messageStore.GetAllMessages()
		data,_ := json.Marshal(msgs)
		if err != nil {
			log.Println("Cant get messages:",err.Error())
		}else{
			conn.WriteMessage(1,data)
		}
		log.Println("Add client. Connection:",len(clients))
		reader(conn,ef.messageStore)
	}
}


// wait message from client
func reader(conn *websocket.Conn,mongoDb MessageStore){
	for {
		msgType,p,err := conn.ReadMessage()

		//remove client from the map, if there is error
		if err !=nil {
			log.Println(err)
			n--
			m := Message{
				UserName: "Server",
				Text:     fmt.Sprintf("1 disconnect"),
			}
			data,_ := json.Marshal(m)
			mu.Lock()
			delete(clients,conn)
			mu.Unlock()
			log.Println("remove client. Connections:",len(clients))
			mu.Lock()
			broadcastMessage(data,1)
			mu.Unlock()
			return
		}

		log.Println(mongoDb.GetAllMessages())
		msg := Message{}
		json.Unmarshal(p,&msg)
		_,err = mongoDb.CreateMessage(&msg)
		if err != nil{
			log.Println("Cannot save the message ",msg)
		}
		log.Println(string(p),msgType)
		mu.Lock()
		broadcastMessage(p,msgType)
		mu.Unlock()
	}
}

// broadcast messages for all clients
func broadcastMessage(msg []byte,msgType int){
	for v,_ := range clients {
		if err := v.WriteMessage(msgType,msg);err!=nil{
			log.Println(err)
			return
		}
	}
}