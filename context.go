package thor

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	ckrouter "github.com/CloudyKit/router"
)

//Context allow us to pass variables between middleware
type Context struct {
	Response ResponseWriter
	Request  *http.Request
	data     map[string]interface{}
	params   ckrouter.Parameter
	handlers []HandlerFunc
	index    int8
	Thor     *Thor
	writer   writer
}

// Next should be used only in the middlewares.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for s := int8(len(c.handlers)); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//ClientIP get ip from RemoteAddr
func (c *Context) ClientIP() string {
	return c.Request.RemoteAddr
}

// Param returns the value of the URL param.
// It is a shortcut for c.params.ByName(key)
//     r.GET("/user/:id", func(c *thor.Context) {
//         // a GET request to /user/onerciller
//         id := c.Param("id") // id == "onerciller"
//     })
func (c *Context) Param(name string) string {
	val := c.params.ByName(name)
	if val == "" {
		val = c.Request.URL.Query().Get(name)
	}
	return val
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.Query("id") == "1234"
// 	   c.Query("name") == "Manu"
// 	   c.Query("value") == ""
// 	   c.Query("wtf") == ""
func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
//     GET /?name=Manu&lastname=
//     c.DefaultQuery("name", "unknown") == "Manu"
//     c.DefaultQuery("id", "none") == "none"
//     c.DefaultQuery("lastname", "none") == ""
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /?name=Manu&lastname=
//     ("Manu", true) == c.GetQuery("name")
//     ("", false) == c.GetQuery("id")
//     ("", true) == c.GetQuery("lastname")
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	if values, ok := c.Request.URL.Query()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
// For example, during a PATCH request to update the user's email:
//     email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
// 	   email=                  -->  ("", true) := GetPostForm("email") // set email to ""
//                             -->  ("", false) := GetPostForm("email") // do nothing with email
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	req := c.Request
	req.ParseForm()
	req.ParseMultipartForm(c.Thor.MaxMultipartMemory)
	if values := req.PostForm[key]; len(values) > 0 {
		return values, true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values, true
		}
	}
	return []string{}, false
}

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	_, fh, err := c.Request.FormFile(name)
	return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.Thor.MaxMultipartMemory)
	return c.Request.MultipartForm, err
}

// SaveUploadedFile uploads the form file to specific dst.
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(out, src)
	return nil
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

// JSON render
func (c *Context) JSON(status int, value interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(value)
	if err != nil {
		panic("couldnt encode response body to json")
	}
	c.Response.Write(json)
}

//Bind decode json to interface{}
func (c *Context) Bind(obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(obj)
	return err
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.data[key]
	return
}

func (c *Context) reset() {
	c.Response = &c.writer
	c.handlers = nil
	c.index = -1
	c.data = nil
}

//reuseContext for resusable context
func (t *Thor) reuseContext(ctx *Context) {
	t.pool.Put(ctx)
}

func (c *Thor) createContext(w http.ResponseWriter, req *http.Request, params ckrouter.Parameter, handlers []HandlerFunc) *Context {
	ctx := c.pool.Get().(*Context)
	ctx.Response = &ctx.writer
	ctx.Request = req
	ctx.data = nil
	ctx.handlers = handlers
	ctx.writer.reset(w)
	ctx.index = -1
	return ctx
}
