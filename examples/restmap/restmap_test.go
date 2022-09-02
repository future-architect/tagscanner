package restmap

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func Test_requestDecoder_ExtractValue(t *testing.T) {
	type args struct {
		tag *RestTag
	}
	tests := []struct {
		name       string
		newRequest func() *http.Request
		args       args
		wantValue  any
		wantErr    bool
	}{
		{
			name: "method",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				return req
			},
			args: args{
				tag: &RestTag{
					Type: MethodField,
				},
			},
			wantValue: "GET",
			wantErr:   false,
		},
		{
			name: "header: found",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set("Accept", "application/json")
				return req
			},
			args: args{
				tag: &RestTag{
					Type: HeaderField,
					Name: "Accept",
				},
			},
			wantValue: "application/json",
			wantErr:   false,
		},
		{
			name: "header: fallback to default",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				return req
			},
			args: args{
				tag: &RestTag{
					Type:    HeaderField,
					Name:    "Accept",
					Default: "application/json",
				},
			},
			wantValue: "application/json",
			wantErr:   false,
		},
		{
			name: "header: not found",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				return req
			},
			args: args{
				tag: &RestTag{
					Type: HeaderField,
					Name: "Accept",
				},
			},
			wantValue: "",
			wantErr:   false,
		},
		{
			name: "cookie: found",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set("Cookie", "cookie1=value1;cookie2=value2")
				return req
			},
			args: args{
				tag: &RestTag{
					Type: CookieField,
					Name: "cookie1",
				},
			},
			wantValue: "value1",
			wantErr:   false,
		},
		{
			name: "cookie: fallback",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set("Cookie", "cookie1=value1;cookie2=value2")
				return req
			},
			args: args{
				tag: &RestTag{
					Type:    CookieField,
					Name:    "cookie3",
					Default: "value3",
				},
			},
			wantValue: "value3",
			wantErr:   false,
		},
		{
			name: "cookie: not found",
			newRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set("Cookie", "cookie1=value1;cookie2=value2")
				return req
			},
			args: args{
				tag: &RestTag{
					Type: CookieField,
					Name: "cookie3",
				},
			},
			wantValue: "",
			wantErr:   false,
		},
		{
			name: "query: found",
			newRequest: func() *http.Request {
				q := url.Values{
					"query1": []string{"value1"},
					"query2": []string{"value2"},
				}
				req, _ := http.NewRequest("GET", "http://example.com?"+q.Encode(), nil)
				return req
			},
			args: args{
				tag: &RestTag{
					Type: QueryField,
					Name: "query1",
				},
			},
			wantValue: "value1",
			wantErr:   false,
		},
		{
			name: "query: fallback to default",
			newRequest: func() *http.Request {
				q := url.Values{
					"query1": []string{"value1"},
					"query2": []string{"value2"},
				}
				req, _ := http.NewRequest("GET", "http://example.com?"+q.Encode(), nil)
				return req
			},
			args: args{
				tag: &RestTag{
					Type:    QueryField,
					Name:    "query3",
					Default: "value3",
				},
			},
			wantValue: "value3",
			wantErr:   false,
		},
		{
			name: "query: not found",
			newRequest: func() *http.Request {
				q := url.Values{
					"query1": []string{"value1"},
					"query2": []string{"value2"},
				}
				req, _ := http.NewRequest("GET", "http://example.com?"+q.Encode(), nil)
				return req
			},
			args: args{
				tag: &RestTag{
					Type: QueryField,
					Name: "query3",
				},
			},
			wantValue: "",
			wantErr:   false,
		},
		{
			name: "body: application/x-www-form-urlencoded",
			newRequest: func() *http.Request {
				q := url.Values{
					"name": []string{"wozozo"},
					"age":  []string{"secret"},
				}
				req, _ := http.NewRequest("POST", "http://example.com", strings.NewReader(q.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			args: args{
				tag: &RestTag{
					Type: BodyField,
					Name: "name",
				},
			},
			wantValue: "wozozo",
			wantErr:   false,
		},
		{
			name: "body: multipart/form-data",
			newRequest: func() *http.Request {
				var b bytes.Buffer
				w := multipart.NewWriter(&b)
				writer, _ := w.CreateFormField("name")
				io.WriteString(writer, "wozozo")
				w.Close()
				req, _ := http.NewRequest("POST", "http://example.com", &b)
				req.Header.Set("Content-Type", w.FormDataContentType())
				return req
			},
			args: args{
				tag: &RestTag{
					Type: BodyField,
					Name: "name",
				},
			},
			wantValue: "wozozo",
			wantErr:   false,
		},
		{
			name: "body: application/json",
			newRequest: func() *http.Request {
				var b bytes.Buffer
				q := map[string]string{
					"name": "wozozo",
					"age":  "secret",
				}
				e := json.NewEncoder(&b)
				e.Encode(&q)
				req, _ := http.NewRequest("POST", "http://example.com", &b)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			args: args{
				tag: &RestTag{
					Type: BodyField,
					Name: "name",
				},
			},
			wantValue: "wozozo",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := requestDecoder{
				req:  tt.newRequest(),
				once: &sync.Once{},
			}
			gotValue, err := d.ExtractValue(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("ExtractValue() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
