package runtimescan_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gitlab.com/osaki-lab/tagscanner/runtimescan"
)

type cmpStrategy struct {
	Path       string
	IgnoreCase bool
}

type compare struct {
	values map[string]interface{}
}

func (c compare) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	if tagStr == "skip" {
		return nil, runtimescan.Skip
	}
	if tagStr == "ignorecase" && elemType.Kind() != reflect.String {
		return nil, fmt.Errorf("the field '%v' is specified as 'ignorecase' but it is not string.", pathStr)
	}
	return &cmpStrategy{
		Path:       pathStr,
		IgnoreCase: tagStr == "ignorecase",
	}, nil
}

func (c *compare) VisitField(tag, value interface{}) (err error) {
	t := tag.(*cmpStrategy)
	if t.IgnoreCase {
		c.values[t.Path] = strings.ToLower(value.(string))
	} else {
		c.values[t.Path] = value
	}
	return nil
}

func (c compare) EnterChild(tag interface{}) (err error) {
	return nil
}

func (c compare) LeaveChild(tag interface{}) (err error) {
	return nil
}

func Compare(s1, s2 interface{}) (bool, []string, error) {
	c1 := compare{
		values: make(map[string]interface{}),
	}
	err := runtimescan.Encode(s1, []string{"cmp"}, &c1)
	if err != nil {
		return false, nil, err
	}
	c2 := compare{
		values: make(map[string]interface{}),
	}
	err = runtimescan.Encode(s2, []string{"cmp"}, &c2)
	if err != nil {
		return false, nil, err
	}
	var unmatch []string
	for k, v1 := range c1.values {
		if c2.values[k] != v1 {
			unmatch = append(unmatch, k)
		}
	}
	for k, v2 := range c2.values {
		v1, ok := c1.values[k]
		if !ok && v1 != v2 {
			unmatch = append(unmatch, k)
		}
	}
	sort.Strings(unmatch)
	return len(unmatch) == 0, unmatch, nil
}

func Example_compare() {
	s1 := struct {
		BothHaveSameValue int
		DifferentValue    int
		DifferentType     int
		Skip1             string `cmp:"skip"`
		IgnoreCase        string `cmp:"ignorecase"`
		OnlyOnS1          bool
	}{
		BothHaveSameValue: 17,
		DifferentValue:    19,
		DifferentType:     11,
		Skip1:             "skip by tag",
		IgnoreCase:        "",
		OnlyOnS1:          true,
	}
	s2 := struct {
		BothHaveSameValue int
		DifferentValue    int
		DifferentType     float64
		Skip2             string `cmp:"skip"`
		IgnoreCase        string `cmp:"ignorecase"`
		OnlyOnS2          bool
	}{
		BothHaveSameValue: 17,
		DifferentValue:    23,
		DifferentType:     1.23,
		Skip2:             "skip by tag",
		IgnoreCase:        "",
		OnlyOnS2:          true,
	}
	equal, diffs, err := Compare(&s1, &s2)
	fmt.Printf("equal: %v\n", equal)
	fmt.Printf("different keys: %v\n", diffs)
	fmt.Printf("err: %v\n", err)
	// Output:
	// equal: false
	// different keys: [DifferentType DifferentValue OnlyOnS1 OnlyOnS2]
	// err: <nil>
}
