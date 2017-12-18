package thor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	passed := false
	testrouter := New()
	testrouter.Group("/api", func(group *RouteGroup) {
		group.GET("/test", func(c *Context) error {
			passed = true
			return nil
		})
	})
	w := performRequest(testrouter, "GET", "/api/test", "")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, passed)

}

func Test_Router(t *testing.T) {
	testRoute("GET", t)
	testRoute("POST", t)
	testRoute("DELETE", t)
	testRoute("PATCH", t)
	testRoute("PUT", t)
	testRoute("OPTIONS", t)
	testRoute("HEAD", t)
}

func testRoute(method string, t *testing.T) {
	passed := false
	m := New()
	switch method {
	case "GET":
		m.GET("", func(ctx *Context) error {
			passed = true
			return nil
		})
	case "POST":
		m.POST("", func(ctx *Context) error {
			passed = true
			return nil
		})
	case "DELETE":
		m.DELETE("", func(ctx *Context) error {
			passed = true
			return nil
		})
	case "PATCH":
		m.PATCH("", func(ctx *Context) error {
			passed = true
			return nil
		})
	case "PUT":
		m.PUT("", func(ctx *Context) error {
			passed = true
			return nil
		})
	case "OPTIONS":
		m.OPTIONS("", func(ctx *Context) error {
			passed = true
			return nil

		})
	case "HEAD":
		m.HEAD("", func(ctx *Context) error {
			passed = true
			return nil
		})
	default:
		panic("unknown method")

	}

	w := performRequest(m, method, "/", "")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, passed)

}
