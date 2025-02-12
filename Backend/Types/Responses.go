package Types

// Response is a struct that represents a response from the server

type Response struct {

	// Message is the message that the server will send to the client

	Message string `json:"Message"`

	Data interface{} `json:"Data"`
}
