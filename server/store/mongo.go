package store

import (
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
)



// init new mongo collection
func NewMongo() (MessageStore,error){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client,err := mongo.Connect(context.TODO(),clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(),nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("chat").Collection("message")

	log.Println("Connected to mongodb")
	return &mongoStore{db: collection},nil
}

// db struct
type mongoStore struct {
	db *mongo.Collection
}

// CreateMessage...
func (m *mongoStore) CreateMessage(msg *Message) (*Message, error){
	insertRes,err := m.db.InsertOne(context.TODO(),msg)
	if err != nil {
		return nil,err
	}
	log.Println("Inserted message:",insertRes.InsertedID)
	return msg,nil
}

// GetAllMessages...
func (m *mongoStore) GetAllMessages() ([]*Message,error){
	findOptions := options.Find()
	var results []*Message
	cur,err := m.db.Find(context.TODO(),bson.D{{}},findOptions)
	if err != nil {
		return nil,err
	}
	for cur.Next(context.TODO()){
		var msg Message
		err := cur.Decode(&msg)
		if err != nil {
			return nil,err
		}
		results = append(results,&msg)
	}
	if err := cur.Err();err != nil {
		return nil,err
	}
	cur.Close(context.TODO())
	return results,nil
}

