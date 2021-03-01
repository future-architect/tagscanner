package staticscan

import (
	"path"
	"reflect"
	"runtime"
	"testing"
)

func testdata(dirname string) string {
	_, filename, _, _ := runtime.Caller(1)
	cwd := path.Dir(filename)
	root := path.Dir(cwd)
	return path.Join(root, "testdata", dirname)
}

func TestScan(t *testing.T) {
	type args struct {
		rootPath string
		tag      string
	}
	tests := []struct {
		name    string
		args    args
		want    []Struct
		wantErr bool
	}{
		{
			name: "simple tags",
			args: args{
				rootPath: testdata("simple"),
				tag: "testtag",
			},
			want: []Struct{
				{
					PackageName: "simple",
					StructName: "SimpleStruct",
					Comment: "SimpleStruct contains multiple fields that have tags\n",
					Fields: []Field{
						{
							Name: "Name",
							Type: "string",
							Tag:  "test1",
						},
						{
							Name: "Age",
							Type: "int",
							Tag:  "test2",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "pointer fields",
			args: args{
				rootPath: testdata("pointer"),
				tag: "testtag",
			},
			want: []Struct{
				{
					PackageName: "pointer",
					StructName: "PointerStruct",
					Comment: "PointerStruct contains multiple pointer fields\n",
					Fields: []Field{
						{
							Name: "Name",
							Type: "*string",
							Tag:  "test1",
						},
						{
							Name: "Age",
							Type: "*int",
							Tag:  "test2",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Scan(tt.args.rootPath, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() got = %v, want %v", got, tt.want)
			}
		})
	}
}
