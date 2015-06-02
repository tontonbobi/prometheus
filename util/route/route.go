package route

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

type param string

func Param(ctx context.Context, p string) string {
	return ctx.Value(param(p)).(string)
}

type Handle func(ctx context.Context, w http.ResponseWriter, req *http.Request)

// handle turns a Handle into httprouter.Handle
func handle(h Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, p := range params {
			ctx = context.WithValue(ctx, param(p.Key), p.Value)
			ctx.Value(param(p))
		}
		h(w, req, ctx)
	}
}

type Router interface {
	http.Handler

	WithPrefix(prefix string) Router
	NewContext() context.Context

	GET(path string, h Handle)
	DELETE(path string, h Handle)
}

type router struct {
	rtr    *httprouter.Router
	prefix string
}

func New() Router {
	return &router{rtr: httprouter.New()}
}

func (r router) WithPrefix(prefix string) Router {
	r.prefix += prefix
	return r
}

func (r *router) GET(path string, h Handle) {
	r.rtr.GET(r.prefix+path, handle(h))
}

func (r *router) DELETE(path string, h Handle) {
	r.rtr.DELETE(r.prefix+path, handle(h))
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}
