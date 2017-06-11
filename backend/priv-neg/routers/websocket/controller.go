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
	render *render.Render
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/ws", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.websocketHandler)),
	)).Methods("GET")

	log.Println("Set up WebSocket controller.")

}

func (c Controller) websocketHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("[websocket] Authenticated %s", middlewares.AuthTokenFromContext(r.Context()))

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
