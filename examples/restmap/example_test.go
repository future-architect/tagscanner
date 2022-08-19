package restmap_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/future-architect/tagscanner/examples/restmap"
)

func ExampleDecode() {
	r := createRequest()

	type Request struct {
		Method  string                `rest:"method"`
		Auth    string                `rest:"header:Authorization"`
		TraceID string                `rest:"cookie:trace-id"`
		Title   string                `rest:"body:title-field"`
		File    multipart.File        `rest:"body:file-field"`
		Header  *multipart.FileHeader `rest:"body:file-field"`
		Ctx     context.Context       `rest:"context"`
	}

	var req Request

	err := restmap.Decode(&req, r)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("method: %s\n", req.Method)
	fmt.Printf("trace-id from cookie: %s\n", req.TraceID)
	fmt.Printf("title from form: %s\n", req.Title)
	fmt.Printf("uploaded file name from form: %s\n", req.Header.Filename)
	content, _ := ioutil.ReadAll(req.File)
	defer req.File.Close()
	fmt.Printf("uploaded file content from form: %s\n", string(content))
	fmt.Printf("value of context: %s\n", req.Ctx.Value("context-key").(string))
	// Output:
	// err: <nil>
	// method: POST
	// trace-id from cookie: 12345
	// title from form: secret file
	// uploaded file name from form: secret.txt
	// uploaded file content from form: file content
	// value of context: context-value
}

func createRequest() *http.Request {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	writer, _ := w.CreateFormField("title-field")
	io.WriteString(writer, "secret file")
	writer, _ = w.CreateFormFile("file-field", "secret.txt")
	io.WriteString(writer, "file content")
	w.Close()

	ctx := context.WithValue(context.Background(), "context-key", "context-value")

	req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "session-key")
	req.Header.Set("Cookie", "trace-id=12345")

	return req
}
