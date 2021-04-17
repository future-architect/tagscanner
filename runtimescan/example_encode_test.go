package runtimescan_test

import (
	"fmt"
	"log"
	"reflect"

	"gitlab.com/osaki-lab/tagscanner/runtimescan"
)

type encoder struct {
	dest map[string]interface{}
}

func (m encoder) ParseTag(name, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	return runtimescan.BasicParseTag(name, tagStr, pathStr, elemType)
}

func (m *encoder) VisitField(tag, value interface{}) (err error) {
	t := tag.(*runtimescan.BasicTag)
	m.dest[t.Tag] = value
	return nil
}

func (m encoder) EnterChild(tag interface{}) (err error) {
	return nil
}

func (m encoder) LeaveChild(tag interface{}) (err error) {
	return nil
}

func Encode(dest map[string]interface{}, src interface{}) error {
	enc := &encoder{
		dest: dest,
	}
	return runtimescan.Encode(src, []string{"map"}, enc)
}

func Example_struct2map() {
	destMap := map[string]interface{}{}
	sampleStruct := struct {
		Int         int      `map:"int"`
		Pointer     *float64 `map:"float"`
		String      string
		NonExisting *bool  `map:"bool"`
		private     string `map:"private"`
	}{
		Int:     13,
		Pointer: &[]float64{3.1415}[0],
		String:  "hello world",
		private: "this should be ignored",
	}
	err := Encode(destMap, &sampleStruct)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("string: %v\n", destMap["string"])
	fmt.Printf("int: %v\n", destMap["int"])
	fmt.Printf("float: %v\n", destMap["float"])
	// Output:
	// string: hello world
	// int: 13
	// float: 3.1415
}
