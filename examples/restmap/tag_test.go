package restmap

import (
	"errors"
	"reflect"
	"testing"

	"github.com/future-architect/tagscanner/runtimescan"
	"github.com/stretchr/testify/assert"
)

func Test_defaultFieldType_Convert(t *testing.T) {
	tests := []struct {
		name string
		t    defaultFieldType
		arg  string
		want string
	}{
		{
			name: "lower",
			t:    lowerCase,
			arg:  "AbcDef",
			want: "abcdef",
		},
		{
			name: "hyphenated-lower-case",
			t:    hyphenatedLowerCase,
			arg:  "AbcDef",
			want: "abc-def",
		},
		{
			name: "hyphenated-pascal-case",
			t:    hyphenatedPascalCase,
			arg:  "AbcDef",
			want: "Abc-Def",
		},
		{
			name: "no field",
			t:    noField,
			arg:  "AbcDef",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Convert(tt.arg); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTag(t *testing.T) {
	type args struct {
		field string
		tag   string
	}
	tests := []struct {
		name         string
		args         args
		wantT        FieldType
		wantName     string
		wantDefault  string
		wantOptional bool
		wantBase     bool
		wantErrType  error
		wantErrMsg   string
	}{
		{
			name: "method",
			args: args{
				field: "Method",
				tag:   "method",
			},
			wantT:       MethodField,
			wantName:    "",
			wantErrType: nil,
		},
		{
			name: "method ng",
			args: args{
				field: "Method",
				tag:   "method:error",
			},
			wantT:       IgnoreField,
			wantName:    "",
			wantErrType: runtimescan.ErrParseTag,
			wantErrMsg:  "method",
		},
		{
			name: "path, default name",
			args: args{
				field: "PathField",
				tag:   "path",
			},
			wantT:       PathField,
			wantName:    "path-field",
			wantErrType: nil,
		},
		{
			name: "path, specified name",
			args: args{
				field: "PathField",
				tag:   "path:user-id",
			},
			wantT:       PathField,
			wantName:    "user-id",
			wantErrType: nil,
		},
		{
			name: "path, can't have default",
			args: args{
				field: "PathField",
				tag:   "path:user-id,default:12345",
			},
			wantErrType: runtimescan.ErrParseTag,
			wantErrMsg:  "can't have default value",
		},
		{
			name: "header, default name",
			args: args{
				field: "ContentType",
				tag:   "header",
			},
			wantT:       HeaderField,
			wantName:    "Content-Type",
			wantErrType: nil,
		},
		{
			name: "header, specified name",
			args: args{
				field: "CType",
				tag:   "header:Content-Type",
			},
			wantT:       HeaderField,
			wantName:    "Content-Type",
			wantErrType: nil,
		},
		{
			name: "header, have default",
			args: args{
				field: "CType",
				tag:   "header:Content-Type,default:application/json",
			},
			wantT:       HeaderField,
			wantName:    "Content-Type",
			wantDefault: "application/json",
			wantErrType: nil,
		},
		{
			name: "header, optional",
			args: args{
				field: "CType",
				tag:   "header:Content-Type,optional",
			},
			wantT:        HeaderField,
			wantName:     "Content-Type",
			wantOptional: true,
			wantErrType:  nil,
		},
		{
			name: "header, base",
			args: args{
				field: "CType",
				tag:   "header:Content-Type,base",
			},
			wantT:       HeaderField,
			wantName:    "Content-Type",
			wantBase:    true,
			wantErrType: nil,
		},
		{
			name: "query, default name",
			args: args{
				field: "PageNumber",
				tag:   "query",
			},
			wantT:       QueryField,
			wantName:    "pagenumber",
			wantErrType: nil,
		},
		{
			name: "query, specified name",
			args: args{
				field: "IPP",
				tag:   "query:item-per-page",
			},
			wantT:       QueryField,
			wantName:    "item-per-page",
			wantErrType: nil,
		},
		{
			name: "query, have default",
			args: args{
				field: "IPP",
				tag:   "query:item-per-page,default:50",
			},
			wantT:       QueryField,
			wantName:    "item-per-page",
			wantDefault: "50",
			wantErrType: nil,
		},
		{
			name: "cookie, default name",
			args: args{
				field: "LastVisited",
				tag:   "cookie",
			},
			wantT:       CookieField,
			wantName:    "last-visited",
			wantErrType: nil,
		},
		{
			name: "cookie, specified name",
			args: args{
				field: "LV",
				tag:   "cookie:last-visited",
			},
			wantT:       CookieField,
			wantName:    "last-visited",
			wantErrType: nil,
		},
		{
			name: "cookie, have default",
			args: args{
				field: "LastVisited",
				tag:   "cookie:last-visited,default:/",
			},
			wantT:       CookieField,
			wantName:    "last-visited",
			wantDefault: "/",
			wantErrType: nil,
		},
		{
			name: "body, default name",
			args: args{
				field: "Name",
				tag:   "body",
			},
			wantT:       BodyField,
			wantName:    "name",
			wantErrType: nil,
		},
		{
			name: "body, specified name",
			args: args{
				field: "UserName",
				tag:   "body:name",
			},
			wantT:       BodyField,
			wantName:    "name",
			wantErrType: nil,
		},
		{
			name: "body, have default",
			args: args{
				field: "Body",
				tag:   "body:name,default:anonymous",
			},
			wantT:       BodyField,
			wantName:    "name",
			wantDefault: "anonymous",
			wantErrType: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRestTag(tt.args.field, tt.args.tag, tt.args.field, reflect.ValueOf("a").Type())
			if got != nil {
				assert.NoError(t, err)
				assert.NoError(t, tt.wantErrType)
				assert.Equal(t, tt.wantT, got.Type)
				assert.Equal(t, tt.wantName, got.Name)
				assert.Equal(t, tt.wantDefault, got.Default)
				assert.Equal(t, tt.wantOptional, got.Optional)
				assert.Equal(t, tt.wantBase, got.Base)
			} else {
				assert.True(t, errors.Is(err, tt.wantErrType))
			}
			if tt.wantErrMsg != "" && err != nil {
				assert.Contains(t, err.Error(), err.Error())
			}
		})
	}
}
