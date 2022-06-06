package staticscan

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testdata(dirname string) string {
	_, filename, _, _ := runtime.Caller(1)
	cwd := filepath.Dir(filename)
	root := filepath.Dir(cwd)
	return filepath.Join(root, "testdata", dirname)
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
				tag:      "testtag",
			},
			want: []Struct{
				{
					PackageName: "simple",
					StructName:  "SimpleStruct",
					Comment:     "SimpleStruct contains multiple fields that have tags\n",
					Fields: []Field{
						{
							Name:     "Name",
							Type:     "string",
							Tag:      "test1",
							FullPath: filepath.Join(testdata("simple"), "simplestruct.go"),
						},
						{
							Name:     "Age",
							Type:     "int",
							Tag:      "test2",
							FullPath: filepath.Join(testdata("simple"), "simplestruct.go"),
						},
					},
					FullPath: filepath.Join(testdata("simple"), "simplestruct.go"),
				},
			},
			wantErr: false,
		},
		{
			name: "pointer fields",
			args: args{
				rootPath: testdata("pointer"),
				tag:      "testtag",
			},
			want: []Struct{
				{
					PackageName: "pointer",
					StructName:  "PointerStruct",
					Comment:     "PointerStruct contains multiple pointer fields\n",
					Fields: []Field{
						{
							Name:     "Name",
							Type:     "*string",
							Tag:      "test1",
							FullPath: filepath.Join(testdata("pointer"), "pointer.go"),
						},
						{
							Name:     "Age",
							Type:     "*int",
							Tag:      "test2",
							FullPath: filepath.Join(testdata("pointer"), "pointer.go"),
						},
					},
					FullPath: filepath.Join(testdata("pointer"), "pointer.go"),
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
			/*for _, s := range got {
				for _, f := range s.Fields {
					f.FullPath = ""
				}
			}*/
			assert.Equal(t, got, tt.want)
		})
	}
}
