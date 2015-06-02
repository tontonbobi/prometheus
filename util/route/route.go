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

type key int

const patKey key = 0

func Pat(ctx context.Context) string {
	return ctx.Value(patKey).(string)
}

type Handle func(w http.ResponseWriter, req *http.Request, ctx context.Context)

// handle turns a Handle into httprouter.Handle
func handle(pat string, h Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, p := range params {
			ctx = context.WithValue(ctx, param(p.Key), p.Value)
		}
		ctx = context.WithValue(ctx, patKey, pat)
		h(w, req, ctx)
	}
}

type Router struct {
	rtr    *httprouter.Router
	prefix string
}

func New() *Router {
	return &Router{rtr: httprouter.New()}
}

func (r *Router) WithPrefix(prefix string) *Router {
	return &Router{rtr: r.rtr, prefix: r.prefix + prefix}
}

func (r *Router) GET(pat string, h Handle) {
	r.rtr.GET(r.prefix+path, handle(pat, h))
}

func (r *Router) DELETE(pat string, h Handle) {
	r.rtr.DELETE(r.prefix+pat, handle(pat, h))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}
