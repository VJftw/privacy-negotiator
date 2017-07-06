package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - WebSocket
type Controller struct {
	logger *log.Logger
	render *render.Render
	redis  *redis.Pool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// NewController - Returns a new controller for websockets.
func NewController(
	webSocketLogger *log.Logger,
	renderer *render.Render,
	redis *redis.Pool,
) *Controller {
	return &Controller{
		logger: webSocketLogger,
		render: renderer,
		redis:  redis,
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

	fbUserID := middlewares.FBUserIDFromContext(r.Context())
	c.logger.Printf("Authenticated %s", fbUserID)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go monitorClose(ws)

	// subscribe to pubsub on redis.
	psc := redis.PubSubConn{Conn: c.redis.Get()}
	psc.Subscribe(fmt.Sprintf("user:%s", fbUserID))

	go redisSubscribe(ws, psc)
}

func redisSubscribe(ws *websocket.Conn, psc redis.PubSubConn) {
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			ws.WriteMessage(websocket.TextMessage, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Printf("error: %v", v)
		}
	}
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
