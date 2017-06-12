package websocket

// Message - for sending a message to be published over a websocket.
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
