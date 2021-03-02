# tagscanner

[![Go Reference](https://pkg.go.dev/badge/gitlab.com/osaki-lab/tagscanner.svg)](https://pkg.go.dev/gitlab.com/osaki-lab/tagscanner)

This package is a helper library for creating a library that interacts with external data using structure tags.

It conceals complex reflection, Go's code analysis, type mapping codes from your code.

It contains three features:

* Extract values from struct (``runtimescan/Encode()``)
* Inject values to struct (``runtimescan/Decode()``)
* Generate code from struct code (``staticscan/Scan()``)

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

#### Write data to other container(``runtimescan/Encode()``)

First, create a structure that satisfies the `` Encoder`` interface. Set the output destination to the field of the structure.

There are also a helper functions.  `` runtimescan.BasicParseTag () `` is a one of them. if you want to put only the field name in the tag, like the basics of `` encoding / json``.

Implementation is completed by setting the value passed to `` VisitField () `` as the output destination.

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

最後に、ユーザー向けのAPIの関数を作ります。

```go
func Encode(dest map[string]interface{}, src interface{}) error {
	enc := &encoder{
		dest: dest,
	}
	return runtimescan.Encode(src, "map", enc)
}
```

#### 外部のデータを構造体の書き込む(``runtimescan/Decode()``)

まずは``Decoder``インタフェースを満たす構造体を作ります。入力元を構造体のフィールドに設定しておきます。

この``ExtractValue()``の返り値が最終的に構造体に書き込まれます。構造体に設定した入力用データから、タグの情報を元に取り出してきて返り値として返すことで、
あとはライブラリがフィールドに値を設定します。

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

こちらも、最後にユーザー向けのAPIの関数を定義します。

```go
func Decode(dest interface{}, src map[string]interface{}) error {
	dec := &decoder{
		src: src,
	}
	return runtimescan.Decode(dest, "map", dec)
}


```

#### 構造体を元にコード生成を行う(``staticscan/Scan()``)

こちらはソースコードの構造体を静的解析します。

#### 実装のヘルパー

いくつか便利関数を提供しています。

* ``runtimescan.BasicParseTag(name, tagStr, pathStr string, elemType reflect.Type)``

  ``ParseTag()``で``encoding/json``のように出力先のキー名（省略時はフィールド名を小文字にしたもの）

* ``runtimescan.Str2PrimitiveValue(v str)``

  1やtrueなどの文字列表現からプリミティブを作成します。タグ中に文字列で書かれたプリミティブをデフォルト値などで使う場合に利用します。

* ``runtimescan.IsPointerOfStruct(i interface{})``, ``runtimescan.IsPointerOfSliceOfStruct(i interface{})``, ``runtimescan.IsPointerOfSliceOfPointerOfStruct(i interface{})``

  ``interface{}``に渡されたポインタ型が``*struct``か``*[]struct``か``*[]*struct``かをそれぞれ判定します。

* ``runtimescan.NewStructInstance(i interface{})``

  引数で渡されたポインタ型を元に、インスタンスを生成して返します。引数は``*struct``か``*[]struct``か``*[]*struct``のいずれかを受け付け、``*struct``を返します。

### 応用例


#### 2つの構造体のインスタンスの比較

``runtimescan/Encode()``をそれぞれのインスタンスごとに呼び、結果を``map``に入れてから比較することで構造体の比較が実現できます。

#### 構造体のコピー

``runtimescan/Encode()``をソースのインスタンスに対して呼び出し、``map``に一時的に値を入れてからそれを元に``runtimescan/Decode()``を呼び出すことで、構造体間でフィールドのコピーが行えます。

## License

Apache2
