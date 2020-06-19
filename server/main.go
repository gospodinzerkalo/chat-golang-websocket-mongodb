package main

import (
	"log"
	"net/http"
	"github.com/gospodinzerkalo/chat/server/store"

)




func main() {

	//connect db
	mongoDb,err := store.NewMongo()
	if err != nil {
		log.Fatal(err)
	}

	// get endpoint factory
	endpoints := store.NewEndpointsFactory(mongoDb)

	//endpoints...
	http.HandleFunc("/",endpoints.WebsocketEndpoint())
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}


