package main

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"fmt"
	"log"
	"os"
)

type Message struct {
	UserName 		string 		`json:"user_name"`
	Text			string 		`json:"text"`
}

func main(){
	conn,_,err := websocket.DefaultDialer.Dial("ws://localhost:9000",nil)
	if err != nil {
		panic(err)
	}
	msgs := []*Message{}
	_,s,_ := conn.ReadMessage()
	if err := json.Unmarshal(s,&msgs);err != nil {
		log.Println(err)
	}
	for _,m := range msgs{
		log.Println(fmt.Sprintf("%v: %v",m.UserName,m.Text))
	}

	go readMessage(conn)
	writeMessage(conn)




	//go readMessage(conn)
	//go writeMessage(conn)
}



func readMessage(conn *websocket.Conn) {
	var m  Message
	for {
		_,p,err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.Unmarshal(p,&m);err!=nil{
			log.Println(err)
			return
		}
		log.Printf("%v: %v",m.UserName,m.Text)
	}
}

func writeMessage(conn *websocket.Conn){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type your username:")

	userName := ""
	nameS := bufio.NewScanner(os.Stdin)
	for nameS.Scan(){
		name := nameS.Text()
		if name == ""{
			fmt.Println("Please, enter your username: ")
			continue
		}
		userName = name
		break
	}
	for scanner.Scan(){
		text := scanner.Text()
		if text == ""{
			continue
		}
		m := &Message{UserName: userName,Text: text}
		data,_ := json.Marshal(m)
		err := conn.WriteMessage(1,data)
		if err != nil {
			log.Println(err)
			break
		}

	}
}