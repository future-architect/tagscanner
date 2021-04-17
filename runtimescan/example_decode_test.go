package runtimescan_test

import (
	"fmt"
	"log"
	"reflect"

	"gitlab.com/osaki-lab/tagscanner/runtimescan"
)

// User should implement runtimescan.Decoder interface
// This instance is created in user code before runtimescan.Decode() function call
type decoder struct {
	src map[string]interface{}
}

func (m decoder) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	return runtimescan.BasicParseTag(name, tagKey, tagStr, pathStr, elemType)
}

func (m *decoder) ExtractValue(tag interface{}) (value interface{}, err error) {
	t := tag.(*runtimescan.BasicTag)
	v, ok := m.src[t.Tag]
	if !ok {
		return nil, runtimescan.Skip
	}
	return v, nil
}

func Decode(dest interface{}, src map[string]interface{}) error {
	dec := &decoder{
		src: src,
	}
	return runtimescan.Decode(dest, []string{"map"}, dec)
}

func Example_map2struct() {
	sampleMap := map[string]interface{}{
		"int":     1,
		"float":   1.1,
		"string":  "hello world",
		"private": "this should be ignored",
	}
	sampleStruct := struct {
		Int         int      `map:"int"`
		Pointer     *float64 `map:"float"`
		String      string
		NonExisting *bool  `map:"bool"`
		private     string `map:"private"`
	}{}
	err := Decode(&sampleStruct, sampleMap)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Field with tag: %v\n", sampleStruct.Int)
	fmt.Printf("Pointer field with tag: %v\n", *(sampleStruct.Pointer))
	fmt.Printf("Field without tag: %v\n", sampleStruct.String)
	fmt.Printf("Field that doesn't exist in source: %v\n", sampleStruct.NonExisting)
	fmt.Printf("Private field is ignored: %v\n", sampleStruct.private)
	// Output:
	// Field with tag: 1
	// Pointer field with tag: 1.1
	// Field without tag: hello world
	// Field that doesn't exist in source: <nil>
	// Private field is ignored:
}
