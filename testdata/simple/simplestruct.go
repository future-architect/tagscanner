package simple

// SimpleStruct contains multiple fields that have tags
type SimpleStruct struct {
	Name string `testtag:"test1"`
	Age  int    `testtag:"test2"`
}
