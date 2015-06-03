package route

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

var (
	mtx   = sync.RWMutex{}
	ctxts = map[*http.Request]context.Context{}
)

// Context returns the context for the request.
func Context(r *http.Request) context.Context {
	mtx.RLock()
	defer mtx.RUnlock()
	return ctxts[r]
}

type param string

func Param(ctx context.Context, p string) string {
	return ctx.Value(param(p)).(string)
}

// handle turns a Handle into httprouter.Handle
func handle(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, p := range params {
			ctx = context.WithValue(ctx, param(p.Key), p.Value)
		}

		mtx.Lock()
		ctxts[r] = ctx
		mtx.Unlock()

		h(w, r)

		mtx.Lock()
		delete(ctxts, r)
		mtx.Unlock()
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

func (r *Router) Get(path string, h http.HandlerFunc) {
	r.rtr.GET(r.prefix+path, handle(h))
}

func (r *Router) Del(path string, h http.HandlerFunc) {
	r.rtr.DELETE(r.prefix+path, handle(h))
}

func (r *Router) Post(path string, h http.HandlerFunc) {
	r.rtr.POST(r.prefix+path, handle(h))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}

func FileServe(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = Param(Context(r), "filepath")
		fs.ServeHTTP(w, r)
	}
}
