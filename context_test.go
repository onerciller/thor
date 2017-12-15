package thor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type ExampleJSON struct {
	Name     string `json:"name"`
	Lastname string `json:"last_name"`
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

func Test_Render(t *testing.T) {
	m := New()
	m.GET("/json", func(ctx *Context) error {
		ctx.JSON(200, &ExampleJSON{Name: "oner", Lastname: "ciller"})
		return nil
	})

	w := performRequest(m, "GET", "/json", "")
	Convey("Json Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(ContentType), ShouldEqual, ContentJSON)
	})

	w = performRequest(m, "GET", "/test", "")
	Convey("Notfound", t, func() {
		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(w.Header().Get(ContentType), ShouldEqual, MIMEPlain)
	})

}

func TestContextQuery(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://google.com/?foo=bar&page=10&id=", nil)

	Convey("Get Query", t, func() {
		value, ok := c.GetQuery("foo")
		So(true, ShouldEqual, ok)
		So("bar", ShouldEqual, value)
	})

	Convey("Default Query", t, func() {
		value := c.DefaultQuery("foo", "none")
		So("bar", ShouldEqual, value)
	})

	// assert.Equal(t, "bar", c.DefaultQuery("foo", "none"))
	// assert.Equal(t, "bar", c.Query("foo"))

	// value, ok = c.GetQuery("page")
	// assert.True(t, ok)
	// assert.Equal(t, "10", value)
	// assert.Equal(t, "10", c.DefaultQuery("page", "0"))
	// assert.Equal(t, "10", c.Query("page"))

	// value, ok = c.GetQuery("id")
	// assert.True(t, ok)
	// assert.Empty(t, value)
	// assert.Empty(t, c.DefaultQuery("id", "nada"))
	// assert.Empty(t, c.Query("id"))

	// value, ok = c.GetQuery("NoKey")
	// assert.False(t, ok)
	// assert.Empty(t, value)
	// assert.Equal(t, "nada", c.DefaultQuery("NoKey", "nada"))
	// assert.Empty(t, c.Query("NoKey"))

	// // postform should not mess
	// value, ok = c.GetPostForm("page")
	// assert.False(t, ok)
	// assert.Empty(t, value)
	// assert.Empty(t, c.PostForm("foo"))
}
