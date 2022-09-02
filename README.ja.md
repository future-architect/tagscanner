# tagscanner

[![Go Reference](https://pkg.go.dev/badge/github.com/future-architect/tagscanner.svg)](https://pkg.go.dev/github.com/future-architect/tagscanner)

このパッケージは構造体のタグを使って外部データとのやりとりをするライブラリを作成するためのヘルパーライブラリです。

リフレクションやGoのコード解析、型のマッピングの処理をカプセル化します。

主に3つの機能があります。

* 構造体のデータを外部に書き出す(``runtimescan.Encode()``)
* 外部のデータを構造体の書き込む(``runtimescan.Decode()``)
* 構造体を元にコード生成を行う(``staticscan.Scan()``)

``runtimescan``パッケージは、実行時に動的に構造体をパースして処理します。
``staticscan``パッケージは、静的解析・コードジェネレータ用です。

## 実行時の処理

用語としては、構造体に書き込む方をデコード、構造体からの読み出しをエンコードと呼んでいます（``encoding/json``と同じ）。

デコードでは``Decoder``インタフェースを、エンコードでは``Encoder``インタフェースを実装します。
インタフェースのインスタンスをそれぞれ、``runtimescan.Decode()``、``runtimescan.Encode()``関数に渡します。

エンコード、デコードの両方で、まずは構造体のフィールドを解析し、上記のインタフェースの``ParseTag()``を呼び出します。
この中でタグの値を分析したりします。このメソッドの返り値は次の処理で利用されます。

デコード処理では``ExtractValue()``が呼ばれます。呼ばれるときには``ParseTag()``の返したタグの分析情報のインスタンスが引数として渡されます。
このメソッドが返した値が、構造体にセットされます。

エンコード処理では``VisitField()``が呼ばれます。こちらも呼ばれるときには``ParseTag()``の返したタグの分析情報のインスタンスが引数として渡されます。
それ以外にフィールドの値も引数として渡されます。

### 基本の使い方

#### 構造体のデータを外部に書き出す(``runtimescan.Encode()``)

まずは``Encoder``インタフェースを満たす構造体を作ります。出力先を構造体のフィールドに設定しておきます。

``encoding/json``の基本のように、タグにはフィールド名のみを入れたい場合は``runtimescan.BasicParseTag()``というヘルパー関数もあります。

``VisitField()``に渡されてくる値を出力先に設定していけば実装完了です。

```go
type encoder struct {
	dest map[string]any
}

func (m encoder) ParseTag(name, tagStr, pathStr string, elemType reflect.Type) (tag any, err error) {
	return runtimescan.BasicParseTag(name, tagStr, pathStr, elemType)
}

func (m *encoder) VisitField(tag, value any) (err error) {
	t := tag.(*runtimescan.BasicTag)
	m.dest[t.Tag] = value
	return nil
}

func (m encoder) EnterChild(tag any) (err error) {
	return nil
}

func (m encoder) LeaveChild(tag any) (err error) {
	return nil
}
```

最後に、ユーザー向けのAPIの関数を作ります。

```go
func Encode(dest map[string]any, src any) error {
	enc := &encoder{
		dest: dest,
	}
	return runtimescan.Encode(src, "map", enc)
}
```

#### 外部のデータを構造体の書き込む(``runtimescan.Decode()``)

まずは``Decoder``インタフェースを満たす構造体を作ります。入力元を構造体のフィールドに設定しておきます。

この``ExtractValue()``の返り値が最終的に構造体に書き込まれます。構造体に設定した入力用データから、タグの情報を元に取り出してきて返り値として返すことで、
あとはライブラリがフィールドに値を設定します。

```go
type decoder struct {
	src map[string]any
}

func (m decoder) ParseTag(name, tagStr, pathStr string, elemType reflect.Type) (tag any, err error) {
	return runtimescan.BasicParseTag(name, tagStr, pathStr, elemType)
}

func (m *decoder) ExtractValue(tag any) (value any, err error) {
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
func Decode(dest any, src map[string]any) error {
	dec := &decoder{
		src: src,
	}
	return runtimescan.Decode(dest, []string{"map"}, dec)
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

* ``runtimescan.IsPointerOfStruct(i any)``, ``runtimescan.IsPointerOfSliceOfStruct(i any)``, ``runtimescan.IsPointerOfSliceOfPointerOfStruct(i any)``

  ``any``に渡されたポインタ型が``*struct``か``*[]struct``か``*[]*struct``かをそれぞれ判定します。``Decode()``を実装する時の型チェックに使います。

* ``runtimescan.NewStructInstance(i any)``

  引数で渡されたポインタ型を元に、インスタンスを生成して返します。引数は``*struct``か``*[]struct``か``*[]*struct``のいずれかを受け付け、``*struct``を返します。

### 応用例

#### 2つの構造体のインスタンスの比較

``runtimescan.Encode()``をそれぞれのインスタンスごとに呼び、結果を``map``に入れてから比較することで構造体の比較が実現できます。

#### 構造体のコピー

``runtimescan.Encode()``をソースのインスタンスに対して呼び出し、``map``に一時的に値を入れてからそれを元に``runtimescan.Decode()``を呼び出すことで、構造体間でフィールドのコピーが行えます。

### サンプル

#### examples/restmap

``runtimescan.Decode``を使い、HTTPリクエストの内容を構造体にマップします。
次のような構造体を作成して、``restmap.Decode()``に``http.Request``とこの構造体のインスタンスをわたします。
``body``のタグはリクエストの``Content-Type``ヘッダーを見て、``application/x-www-form-urlencoded``, ``multipart/form-data``, ``application/json``のいづれかであればパースしてデータを読み込みます。
ファイルのアップロードにも対応します。``rest:"path:param"`で、パスの一部から文字列を取り出しますが、今のところchi routerにのみ対応しています。

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

``runtimescan.Decode``を使い、バイナリをロードします。タグに記述されたルールに従い、バイト列を読み込みアサインします。

* ``数値``でビット数かバイト数を指定します
* ``<< >>``で即値を指定します。文字列か数値が指定可能です。数値の場合は``/``でビット数かバイト数を指定します

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

## 静的なタグの取得

ソースコードを静的スキャンして構造体情報を取り出します。タグを元にしたコード生成のための機能です。

``staticscan.Scan(rootPath, tagName: string) ([]staticscan.Struct, error)``

## ライセンス

Apache2
