package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - WebSocket
type Controller struct {
	logger *log.Logger
	render *render.Render
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// NewController - Returns a new controller for websockets.
func NewController(webSocketLogger *log.Logger, renderer *render.Render) *Controller {
	return &Controller{
		logger: webSocketLogger,
		render: renderer,
	}
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/ws", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.websocketHandler)),
	)).Methods("GET")

	log.Println("Set up WebSocket controller.")

}

func (c Controller) websocketHandler(w http.ResponseWriter, r *http.Request) {

	c.logger.Printf("Authenticated %s", middlewares.FBUserIDFromContext(r.Context()))

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go monitorClose(ws)

	// subscribe to pubsub on redis.
}

func monitorClose(ws *websocket.Conn) {
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing WebSocket")
			ws.Close()
			return
		}
	}
}
