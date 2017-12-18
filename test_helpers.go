package thor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
)

// CreateTestContext returns a fresh engine and context for testing purposes
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Thor) {
	r = New()
	c = r.AllocateContext()
	c.reset()
	c.writer.reset(w)
	return
}

func performRequest(r http.Handler, method, path string, postData string) *httptest.ResponseRecorder {

	req, _ := http.NewRequest(method, path, nil)

	if strings.ToLower(method) == "post" {
		data, _ := url.ParseQuery(postData)
		req, _ = http.NewRequest(method, path, bytes.NewBufferString(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	} else if strings.ToLower(method) == "post|json" {
		req, _ = http.NewRequest("POST", path, bytes.NewBufferString(postData))
		req.Header.Add("Content-Type", "application/json;")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
