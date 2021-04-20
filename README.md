# tagscanner

[![Go Reference](https://pkg.go.dev/badge/gitlab.com/osaki-lab/tagscanner.svg)](https://pkg.go.dev/gitlab.com/osaki-lab/tagscanner)

This package is a helper library for creating a library that interacts with external data using structure tags.

It conceals complex reflection, Go's code analysis, type mapping codes from your code.

It contains three features:

* Extract values from struct (``runtimescan.Encode()``)
* Inject values to struct (``runtimescan.Decode()``)
* Generate code from struct code (``staticscan.Scan()``)

``runtimescan`` package dynamically analyses struct instances and works. 
``staticscan`` package works for static analysis and code generation.

## Basic Behavior

I define the word in this library (it is as same as ``encoding/json``):

Encode: extract values from struct
Decode: inject values to struct

For decoding, user implements type that satisfies ``Decoder`` interface. For encoding, user does ``Encoder`` interface.

``runtimescan.Decode()`` and ``runtimescan.Encode`` will receive these instance. In both cases, this library calls ``ParseTag()`` method.
User's logic analyses tag string and returns the information. This information will be passed to the next methods.

For decoding, ``ExtractValue()`` will be called. It receives tag information instance that ``ParseTag()`` returns.
Return value of``ExtractValue()`` will be passed to instance of struct.

For encoding, ``VisitField()`` will be called. It receives tag information instance that ``ParseTag()`` returns.
In addition to this, ``VisitField()`` will be received the value that is extracted from struct instance.

### Basic Usages

#### Write data from struct's instance to other container(``runtimescan.Encode()``)

First, create a structure that satisfies the ``runtimescan.Encoder`` interface. Set the output destination to the field of the structure.

There are some helper functions.  ``runtimescan.BasicParseTag()`` is a one of them. if you want to put only the field name in the tag, like ``encoding/json``.

Then implement ``VisitField()`` that receives field value and stores to destination object (in this case, ``dest``).

```go
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
```

At last, create an entry point function.

```go
func Encode(dest map[string]interface{}, src interface{}) error {
	enc := &encoder{
		dest: dest,
	}
	return runtimescan.Encode(src, "map", enc)
}
```

#### Write data to struct(``runtimescan/Decode()``)

First, create a structure that satisfies the ``runtimescan.Decoder`` interface.

Implement ``ExtractValue()`` method. The result value of this method will be passed to target struct.

```go
type decoder struct {
	src map[string]interface{}
}

func (m decoder) ParseTag(name, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	return runtimescan.BasicParseTag(name, tagStr, pathStr, elemType)
}

func (m *decoder) ExtractValue(tag interface{}) (value interface{}, err error) {
	t := tag.(*runtimescan.BasicTag)
	v, ok := m.src[t.Tag]
	if !ok {
		return nil, runtimescan.Skip
	}
	return v, nil
}
```

At last, create an entry point function.

```go
func Decode(dest interface{}, src map[string]interface{}) error {
	dec := &decoder{
		src: src,
	}
	return runtimescan.Decode(dest, []string{"map"}, dec)
}


```

#### Generation code from structs' tag fields(``staticscan.Scan()``)

This package provides functions to analyze and generate codes(``staticscan.Scan()`)

#### Utility functions

This package contains several utility functions.

* ``runtimescan.BasicParseTag(name, tagStr, pathStr string, elemType reflect.Type)``

  This is the simplest implementation of ``ParseTag()``. It returns tag value or lower field name if the field doesn't have tag.

* ``runtimescan.Str2PrimitiveValue(v str)``

  It generates primitive from string like "1", "true". It is for creating primitive from string in tag.

* ``runtimescan.IsPointerOfStruct(i interface{})``, ``runtimescan.IsPointerOfSliceOfStruct(i interface{})``, ``runtimescan.IsPointerOfSliceOfPointerOfStruct(i interface{})``

  Check the pointer passed as ``interface{}`` is ``*struct`` or ``*[]struct`` or ``*[]*struct``.
  This is useful to check type definition in ``Decode`` function.

* ``runtimescan.NewStructInstance(i interface{})``

  Generate a new instance based on passed type. Whether the input is ``*struct`` or ``*[]struct`` or ``*[]*struct``, it returns ``*struct``.

### Advanced usage samples

#### Compare two structure

Call ``runtimescan.Encode()`` for each instance and stores the structs' fields into ``map``. Then compare the result in ``map``.

#### Copy structure

Call ``runtimescan.Encode()`` for source instance and stores the struct's fields into ``map``. Then call ``runtimescan.Decode()``.

### Examples

#### examples/restmap

This maps HTTP request to struct instance by using ``runtimescan.Decode``.
Prepare struct like the following and pass its instance and ``http.Request`` to ``restmpa.Decode()``.
``body`` tag detects body format by checking ``Content-Type`` header and accepts ``application/x-www-form-urlencoded``, ``multipart/form-data``, ``application/json`` format.
It also support file uploading. ``rest:"path:param"`` can extract one of requested path when you uses chi router.

```go
type Request struct {
	Method  string                `rest:"method"`
	Auth    string                `rest:"header:Authorization"`
	TraceID string                `rest:"cookie:trace-id"`
	Title   string                `rest:"body:title-field"`
	File    multipart.File        `rest:"body:file-field"`
	Header  *multipart.FileHeader `rest:"body:file-field"`
	Ctx     context.Context       `rest:"context"`
}
```

#### examples/binarypatternmatch

Loads binary data by using ``runtimescan.Decode``. It assigns byte arrays into struct's instance.

* ``num`` specifies length (bits or bytes)
* ``<< >>`` specifies literal. You can use string or number. The number can has length by using ``/``.

```go
type Image struct {
	Header     string `bytes:"<<HEAD>>"`
	Height     int32  `bytes:"4"`
    Width      int32  `bytes:"4"`
	ReadOnly   bool   `bits:"1"`
	_          byte   `bits:"3"`
	ColorType  byte   `bits:"4"'`
	CheckDigit byte   `bits:"<<0x5/6>>"`
}

func main() {
	var image Image
	f, _ := os.Open("imagefile")
    err := binarypatternmatch.Decode(&image, f)
}

```

## Search tags statically

Extrude struct's information by parsing codes statically.
This is for code generation.

``staticscan.Scan(rootPath, tagName: string) ([]staticscan.Struct, error)``


## License

Apache2
