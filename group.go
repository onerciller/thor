package thor

import (
	"net/http"

	ckrouter "github.com/CloudyKit/router"
)

const (
	//GET Request type
	GET = "GET"
	//POST request type
	POST = "POST"
	//PATCH request type
	PATCH = "PATCH"
	//PUT request type
	PUT = "PUT"
	//OPTIONS request type
	OPTIONS = "OPTIONS"
	//CONNECT request type
	CONNECT = "CONNECT"
	//HEAD request type
	HEAD = "HEAD"
	//DELETE request type
	DELETE = "DELETE"
)

//Group struct
type Group struct {
	Handlers     []HandlerFunc
	absolutePath string
	thor         *Thor
}

// Use method is  adds middlewares to the group
func (r *Group) Use(middlewares ...HandlerFunc) {
	r.Handlers = append(r.Handlers, middlewares...)
}

// Group Creates a new router group.
func (r *Group) Group(relativePath string, fn func(*Group), handlers ...HandlerFunc) *Group {
	router := &Group{
		Handlers:     r.combineHandlers(handlers),
		absolutePath: r.calculateAbsolutePath(relativePath),
		thor:         r.thor,
	}
	fn(router)
	return router
}

func (r *Group) calculateAbsolutePath(relativePath string) string {
	return joinPaths(r.absolutePath, relativePath)
}

//Handle method
func (r *Group) Handle(httpMethod, relativePath string, handlers []HandlerFunc) {
	absolutePath := r.calculateAbsolutePath(relativePath)
	handlers = r.combineHandlers(handlers)
	r.thor.router.AddRoute(httpMethod, absolutePath, func(w http.ResponseWriter, req *http.Request, params ckrouter.Parameter) {
		ctx := r.thor.createContext(w, req, params, handlers)
		ctx.Next()
		r.thor.reuseContext(ctx)
	})
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Group) GET(relativePath string, handlers ...HandlerFunc) {
	r.Handle(GET, relativePath, handlers)
}

//POST handle POST method
func (r *Group) POST(path string, handlers ...HandlerFunc) {
	r.Handle(POST, path, handlers)
}

//PATCH handle PATCH method
func (r *Group) PATCH(path string, handlers ...HandlerFunc) {
	r.Handle(PATCH, path, handlers)
}

//PUT handle PUT method
func (r *Group) PUT(path string, handlers ...HandlerFunc) {
	r.Handle(PUT, path, handlers)
}

//DELETE handle DELETE method
func (r *Group) DELETE(path string, handlers ...HandlerFunc) {
	r.Handle(DELETE, path, handlers)
}

//HEAD handle HEAD method
func (r *Group) HEAD(path string, handlers ...HandlerFunc) {
	r.Handle(HEAD, path, handlers)
}

//OPTIONS handle OPTIONS method
func (r *Group) OPTIONS(path string, handlers ...HandlerFunc) {
	r.Handle(OPTIONS, path, handlers)
}

//CONNECT handle OPTIONS method
func (r *Group) CONNECT(path string, handlers ...HandlerFunc) {
	r.Handle(CONNECT, path, handlers)
}

func (r *Group) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(r.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, r.Handlers...)
	return append(mergedHandlers, handlers...)
}
