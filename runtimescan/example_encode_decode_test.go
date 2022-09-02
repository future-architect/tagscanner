package runtimescan_test

import (
	"fmt"
	"log"
	"reflect"

	"github.com/future-architect/tagscanner/runtimescan"
)

type cpyStrategy struct {
	Path string
}

type cpy struct {
	values map[string]any
}

func (c *cpy) VisitField(tag, value any) (err error) {
	t := tag.(*cpyStrategy)
	c.values[t.Path] = value
	return nil
}

func (c cpy) EnterChild(tag any) (err error) {
	return nil
}

func (c cpy) LeaveChild(tag any) (err error) {
	return nil
}

func (c cpy) ExtractValue(tag any) (value any, err error) {
	t := tag.(*cpyStrategy)
	if v, ok := c.values[t.Path]; ok {
		return v, nil
	}
	return nil, runtimescan.Skip
}

func (c cpy) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag any, err error) {
	if tagStr == "skip" {
		return nil, runtimescan.Skip
	}
	return &cpyStrategy{
		Path: pathStr,
	}, nil

}

func Copy(dest, src any) error {
	c := &cpy{
		values: make(map[string]any),
	}
	err := runtimescan.Encode(src, []string{"copy"}, c)
	if err != nil {
		return err
	}
	return runtimescan.Decode(dest, []string{"copy"}, c)
}

type Struct struct {
	Value  string
	Ignore string `copy:"skip"`
}

func Example_copy() {

	src := Struct{
		Value:  "copy from source",
		Ignore: "this value should be ignored",
	}
	dest := Struct{}
	err := Copy(&dest, &src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Value: %s\n", dest.Value)
	fmt.Printf("Ignore: %s\n", dest.Ignore)
	// Output:
	// Value: copy from source
	// Ignore:
}
