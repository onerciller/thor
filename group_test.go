package thor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	passed := false
	thor := New()
	g1 := thor.Group("/api")
	g1.GET("/test", func(c *Context) error {
		passed = true
		return nil
	})

	g2 := g1.Group("/123")

	g2.GET("/sub", func(c *Context) error {
		return nil
	})

	w := performRequest(thor, "GET", "/api/test", "")

	assert.Equal(t, "/api/test", g1.GetFullPath())
	assert.Equal(t, "/api/test/123/sub", g2.GetFullPath())
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
