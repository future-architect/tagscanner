package restmap

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/future-architect/tagscanner/runtimescan"

	"github.com/go-chi/chi/v5"
)

type formFile struct {
	header *multipart.FileHeader
	file   multipart.File
}

type bodyType int

const (
	bodyUnread bodyType = iota
	bodyUrlEncoding
	bodyMultipart
	bodyJson
	bodyUnknown
)

type requestDecoder struct {
	req        *http.Request
	maxMemory  int
	once       *sync.Once
	parseError error
	bodyType   bodyType
	multipart  *multipart.Form
	json       map[string]any
}

func (d requestDecoder) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag any, err error) {
	return ParseRestTag(name, tagStr, pathStr, elemType)
}

func (d *requestDecoder) initBody() {
	ct := d.req.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		d.bodyType = bodyJson
		dec := json.NewDecoder(d.req.Body)
		defer d.req.Body.Close()
		d.parseError = dec.Decode(&d.json)
	} else if ct == "application/x-www-form-urlencoded" {
		d.bodyType = bodyUrlEncoding
		d.parseError = d.req.ParseForm()
	} else if strings.HasPrefix(ct, "multipart/form-data") {
		d.bodyType = bodyMultipart
		maxMemory := d.maxMemory
		if maxMemory == 0 {
			maxMemory = DefaultMaxMemory
		}
		reader, err := d.req.MultipartReader()
		if err != nil {
			d.parseError = err
		} else {
			form, err := reader.ReadForm(int64(maxMemory))
			d.parseError = err
			if err == nil {
				d.multipart = form
			}
		}
	} else {
		d.bodyType = bodyUnknown
	}
}

func (d *requestDecoder) ExtractValue(tagInstance any) (value any, err error) {
	t := tagInstance.(*RestTag)
	switch t.Type {
	case MethodField:
		return d.req.Method, nil
	case PathField:
		return chi.URLParam(d.req, t.Name), nil
	case HeaderField:
		v := d.req.Header.Get(t.Name)
		if v == "" {
			return t.Default, nil
		}
		return v, nil
	case CookieField:
		c, err := d.req.Cookie(t.Name)
		if err == http.ErrNoCookie {
			return t.Default, nil
		} else if err != nil {
			return nil, err
		}
		return c.Value, nil
	case QueryField:
		v := d.req.URL.Query().Get(t.Name)
		if v == "" {
			return t.Default, nil
		}
		// todo: slice support
		return v, nil
	case BodyField:
		d.once.Do(d.initBody)
		if d.parseError != nil {
			return nil, d.parseError
		}
		switch d.bodyType {
		case bodyJson:
			v, ok := d.json[t.Name]
			if ok {
				return v, nil
			}
			return "", nil
		case bodyUrlEncoding:
			return d.req.FormValue(t.Name), nil
		case bodyMultipart:
			var et string
			if t.EType != nil { // for testing
				et = t.EType.String()
			}
			if et == "multipart.File" {
				v, ok := d.multipart.File[t.Name]
				// todo: slice support
				if ok && len(v) > 0 {
					return v[0].Open()
				} else {
					return nil, runtimescan.Skip
				}
			} else if et == "multipart.FileHeader" {
				v, ok := d.multipart.File[t.Name]
				// todo: slice support
				if ok && len(v) > 0 {
					return v[0], nil
				} else {
					return nil, runtimescan.Skip
				}
			} else {
				v, ok := d.multipart.Value[t.Name]
				// todo: slice support
				if !ok {
					return nil, runtimescan.Skip
				}
				return strings.Join(v, ","), nil
			}
		case bodyUnknown:
			return "", nil
		case bodyUnread:
			panic("bug")
		default:
			panic("bug")
		}
	case ContextField:
		return d.req.Context(), nil
	case IgnoreField:
		return nil, runtimescan.Skip
	}
	return nil, runtimescan.Skip
}

var DefaultMaxMemory = 32 << 20 // 32 MB as same as http.Request

var _ runtimescan.Decoder = &requestDecoder{}

func Decode(dest any, r *http.Request) error {
	decoder := &requestDecoder{
		req:      r,
		once:     &sync.Once{},
		bodyType: bodyUnread,
	}
	return runtimescan.Decode(dest, []string{"rest"}, decoder)
}
