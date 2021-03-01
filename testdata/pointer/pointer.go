package pointer

// PointerStruct contains multiple pointer fields
type PointerStruct struct {
	Name *string `testtag:"test1"`
	Age  *int    `testtag:"test2"`
}

