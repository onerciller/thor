package thor

import "net/http"

// CreateTestContext returns a fresh engine and context for testing purposes
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Thor) {
	r = New()
	c = r.AllocateContext()
	c.reset()
	c.writer.reset(w)
	return
}
