package routers

// MuxRouter - The application router
import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// MuxRouter - A GorillaMux router.
type MuxRouter struct {
	Router  *mux.Router
	Render  *render.Render
	Negroni *negroni.Negroni
}

// Routable - Controllers should implement this.
type Routable interface {
	Setup(*mux.Router, *render.Render)
}

// NewMuxRouter - Sets up and returns a new MuxRouter with the given controllers.
func NewMuxRouter(controllers []Routable, logging bool) *MuxRouter {
	muxRouter := &MuxRouter{}

	muxRouter.Render = render.New()
	muxRouter.Router = mux.NewRouter()

	muxRouter.Negroni = negroni.Classic()
	muxRouter.Negroni.Use(cors.New(cors.Options{
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
	}))

	// muxRouter.Handler = n

	for _, controller := range controllers {
		controller.Setup(muxRouter.Router, muxRouter.Render)
	}

	muxRouter.Negroni.UseHandler(muxRouter.Router)

	return muxRouter
}
