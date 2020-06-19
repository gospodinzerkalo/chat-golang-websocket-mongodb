package store

// model message
type Message struct {
	UserName string `json:"user_name"`
	Text	 string `json:"text"`
}


// interfaces for messages
type MessageStore interface {
	CreateMessage(msg *Message) (*Message,error)
	GetAllMessages() ([]*Message,error)
}
