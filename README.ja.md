# tagscanner

このパッケージは構造体のタグを使って外部データとのやりとりをするライブラリを作成するためのヘルパーライブラリです。

リフレクションやGoのコード解析、型のマッピングの処理をカプセル化します。

主に3つの機能があります。

* 構造体のデータを外部に書き出す(``runtimescan/Encode()`)
* 外部のデータを構造体の書き込む(``runtimescan/Decode()`)
* 構造体を元にコード生成を行う(``staticscan/Scan()`)

``runtimescan``パッケージは、実行時に動的に構造体をパースして処理します。
``staticscan``パッケージは、静的解析・コードジェネレータ用です。

## 実行時の処理

用語としては、構造体に書き込む方をデコード、構造体からの読み出しをエンコードと呼んでいます（``encoding/json``と同じ）。

デコードでは``Decoder``インタフェースを、エンコードでは``Encoder``インタフェースを実装します。
インタフェースのインスタンスをそれぞれ、``runtimescan.Decode()``、``runtimescan.Encode()``関数に渡します。

エンコード、デコードの両方で、まずは構造体のフィールドを解析し、上記のインタフェースの``ParseTag()``を呼び出します。
この中でタグの値を分析したりします。このメソッドの返り値は次の処理で利用されます。

デコード処理では``ExtractValue()``が、エンコード処理では``VisitField()``が呼ばれます。

デコードでは``ParseTag()``の返したインスタンスが引数となって帰ってきます。この関数の返り値が構造体にセットされます。

エンコードでは``ExtraceValue()``には``ParseTag()``の返したインスタンスと値が渡ってきます。

前者は

### 基本の使い方

#### 構造体のデータを外部に書き出す(``runtimescan/Encode()``)

まずは``Encoder``インタフェースを満たす構造体を作ります。出力先を構造体のフィールドに設定しておきます。

``encoding/json``の基本のように、タグにはフィールド名のみを入れたい場合は``runtimescan.BasicParseTag()``というヘルパー関数もあります。

``VisitField()``に渡されてくる値を出力先に設定していけば実装完了です。

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

#### 実装のヘルパー



### 応用例


#### 2つの構造体のインスタンスの比較

``runtimescan/Encode()``をそれぞれのインスタンスごとに呼び、結果を``map``に入れてから比較することで構造体の比較が実現できます。

#### 構造体のコピー



## ライセンス

Apache2
