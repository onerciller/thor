package thor

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

const (
	appName                = "THOR"
	defaultMultipartMemory = 32 << 20 // 32 MB

)

//HandlerFunc is middleware func
type HandlerFunc func(*Context) error

//Thor inital struct
type Thor struct {
	*RouteGroup
	router             *httprouter.Router
	pool               sync.Pool
	AppName            string
	MaxMultipartMemory int64
}

// New returns a new blank Thor
func New() *Thor {
	thor := &Thor{}
	thor.AppName = appName
	thor.RouteGroup = &RouteGroup{
		prefix: "/",
		thor:   thor,
	}

	thor.router = httprouter.New()
	thor.MaxMultipartMemory = defaultMultipartMemory
	thor.pool.New = func() interface{} {
		return thor.AllocateContext()
	}
	return thor
}

//AllocateContext is reusable context using pool
func (t *Thor) AllocateContext() *Context {
	return &Context{Thor: t}
}

//Use method for appending middleware
func (t *Thor) Use(middlewares ...HandlerFunc) {
	t.RouteGroup.Use(middlewares...)
}

//ServeHTTP makes the router implement the http.Handler interface.
func (t *Thor) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t.router.ServeHTTP(res, req)

}

// Run run the http server.
func (t *Thor) Run(addr string) error {
	fmt.Printf("%s\n\n", banner)
	fmt.Printf("[%s] Listening and serving HTTP on %s \n\n", t.AppName, addr)

	if err := http.ListenAndServe(addr, t); err != nil {
		return err
	}
	return nil
}
