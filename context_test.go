package thor

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExampleJSON struct {
	Name     string `json:"name"`
	Lastname string `json:"last_name"`
}

func TestContextQuery(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://google.com/?foo=bar&page=10&id=", nil)
	assert.Equal(t, "bar", c.Query("foo"))
}

func TestContextQueryArray(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://google.com/?id=3&array[]=array1&array[]=array2", nil)
	values, ok := c.GetQueryArray("array[]")
	assert.True(t, true, ok)
	assert.Equal(t, "array1", values[0])
}

func TestPostForm(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString("page=test")
	c.Request, _ = http.NewRequest("POST", "http://google.com?test=bar", body)
	c.Request.Header.Add("Content-Type", MIMEPOSTForm)

	values, ok := c.GetPostForm("page")

	assert.True(t, true, ok)
	assert.Equal(t, "test", values)

}

func TestContextFormFile(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		w.Write([]byte("test"))
	}
	mw.Close()
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.FormFile("file")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", f.Filename)
	}

	assert.NoError(t, c.SaveUploadedFile(f, "test"))
}

func TestContextMultipartForm(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.WriteField("foo", "bar")
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		w.Write([]byte("test"))
	}
	mw.Close()
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.MultipartForm()
	if assert.NoError(t, err) {
		assert.NotNil(t, f)
	}

	assert.NoError(t, c.SaveUploadedFile(f.File["file"][0], "test"))
}

func TestSaveUploadedOpenFailed(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.Close()

	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	f := &multipart.FileHeader{
		Filename: "file",
	}
	assert.Error(t, c.SaveUploadedFile(f, "test"))
}

func TestSaveUploadedCreateFailed(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		w.Write([]byte("test"))
	}
	mw.Close()
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.FormFile("file")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", f.Filename)
	}

	assert.Error(t, c.SaveUploadedFile(f, "/"))
}

// TestContextSetGet tests that a parameter is set correctly on the
// current context and can be retrieved using Get.
func TestContextSetGet(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("foo", "bar")

	value, err := c.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, err)

	value, err = c.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, err)
}
