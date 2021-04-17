package runtimescan

import (
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyVisitor struct {
}

func (dv dummyVisitor) ParseTag(name, tag, pathStr string, eType reflect.Type) (interface{}, error) {
	if tag == "" {
		return nil, Skip
	}
	return nil, nil
}

func Test_visitor_parse(t *testing.T) {
	type args struct {
		vi Parser
		e  func() interface{}
	}
	tests := []struct {
		name             string
		args             args
		wantFieldIndexes []int
		wantFieldOps     []visitOpType
		wantError        bool
		wantFieldCount   int
	}{
		{
			name: "single tag",
			args: args{
				vi: &dummyVisitor{},
				e: func() interface{} {
					type S struct {
						I int `rest:"i"`
					}
					return &S{}
				},
			},
			wantFieldIndexes: []int{
				0,
			},
			wantFieldOps: []visitOpType{
				visitFieldOp,
			},
			wantError:      false,
			wantFieldCount: 1,
		},
		{
			name: "single tags",
			args: args{
				vi: &dummyVisitor{},
				e: func() interface{} {
					type S struct {
						I int `rest:"i"`
						S int `rest:"s"`
					}
					return &S{}
				},
			},
			wantFieldIndexes: []int{
				0,
				1,
			},
			wantFieldOps: []visitOpType{
				visitFieldOp,
				visitFieldOp,
			},
			wantError:      false,
			wantFieldCount: 2,
		},
		{
			name: "embed (1)",
			args: args{
				vi: &dummyVisitor{},
				e: func() interface{} {
					type E struct {
						I int `rest:"i"`
					}
					type S struct {
						E
						S int `rest:"s"`
					}
					return &S{}
				},
			},
			wantFieldIndexes: []int{
				0,
				0,
				-1,
				1,
			},
			wantFieldOps: []visitOpType{
				visitChildOp,
				visitFieldOp,
				leaveChildOp,
				visitFieldOp,
			},
			wantError:      false,
			wantFieldCount: 2,
		},
		{
			name: "embed (2)",
			args: args{
				vi: &dummyVisitor{},
				e: func() interface{} {
					type E struct {
						I int `rest:"i"`
					}
					type S struct {
						S int `rest:"s"`
						E
					}
					return &S{}
				},
			},
			wantFieldIndexes: []int{
				0,
				1,
				0,
				-1,
			},
			wantFieldOps: []visitOpType{
				visitFieldOp,
				visitChildOp,
				visitFieldOp,
				leaveChildOp,
			},
			wantError:      false,
			wantFieldCount: 2,
		},
		{
			name: "struct",
			args: args{
				vi: &dummyVisitor{},
				e: func() interface{} {
					type C struct {
						I int `rest:"i"`
					}
					type S struct {
						C C
						S int `rest:"s"`
					}
					return &S{}
				},
			},
			wantFieldIndexes: []int{
				0,
				0,
				-1,
				1,
			},
			wantFieldOps: []visitOpType{
				visitChildOp,
				visitFieldOp,
				leaveChildOp,
				visitFieldOp,
			},
			wantError:      false,
			wantFieldCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &parser{}
			val := reflect.ValueOf(tt.args.e())
			err := v.parse(tt.args.vi, []string{"rest"}, val.Elem().Type())
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFieldIndexes, v.fieldIndexes)
				assert.Equal(t, tt.wantFieldOps, v.fieldOps)
				validFieldCount := 0
				for _, f := range v.fields {
					if f != nil {
						validFieldCount++
					}
				}
				assert.Equal(t, tt.wantFieldCount, validFieldCount)
			}
		})
	}
}

type TestStruct struct {
	FileHeader *multipart.FileHeader `rest:"file"`
	FileFile   multipart.File        `rest:"file"`
}

type dummyStructVisitor struct {
}

func (dv dummyStructVisitor) ParseTag(name, tag, pathStr string) (interface{}, error) {
	if tag == "" {
		return nil, Skip
	}
	return nil, nil
}

/*func TestVisitStructInterface(t *testing.T) {
	v := &parser{}
}*/

/*func TestMapRequestGet(t *testing.T) {
	var result TestStruct
	r := chi.NewRouter()
	r.Get("/book/{isbn}/{page}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		structtag.Decode(&result, r)
	})

	s := httptest.NewServer(r)
	u, _ := url.Parse(s.URL)
	u.Path = "/book/9784873119038/200"
	q := url.Values{
		"line": []string{"10"},
		"highlight": []string{"HTTP", "Method"},
	}
	u.RawQuery = q.Encode()
	t.Log(u.String())
	res, err := http.Get(u.String())
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	assert.Equal(t, "GET", result.Method)
}*/
