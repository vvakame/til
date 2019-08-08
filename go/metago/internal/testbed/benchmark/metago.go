//+build metago

package benchmark

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/vvakame/til/go/metago"
)

type FooMetago struct {
	ID        int64
	Kind      string
	Name      string `json:"nickname"`
	Age       int
	CreatedAt time.Time
}

var propertyNameCache map[string]string

func (obj *FooMetago) MarshalJSON() ([]byte, error) {
	// 実装方法メモ
	//   * bytes.Buffer 作戦 string → []byte のキャストでメモリコピー発生するのを避けるのがだるい
	//     * sync.Pool とかで使い回すとちょっと性能向上
	//     * buf.Reset とかでメモリの再確保を切り詰められる
	//   * strings.Builder
	//     * bytes.Buffer と Reset の実装がことなるので使いまわしに向いてない
	//   * []byte + append(buf, foo...) でひたすらがんばる
	//     * 一番速い
	//     * strconv.AppendQuote とか案外活用できるものがたくさんある
	//     * 読みづらい
	// json.Marshalメモ
	//   * json.Marshal(obj) と obj.MarshalJSON() の結果は等価ではない
	//     * compactとか呼ばれるので
	//     * 前者は結果の []byte が本当にJSONか検証されて出てくる
	//   * 自前 MarshalJSON+json.Marshal は実装無し素json.Marshalに速度的には勝つのは難しい
	//     * 上記の理由により…

	var buf strings.Builder
	if propertyNameCache == nil {
		propertyNameCache = make(map[string]string)
	}

	buf.WriteString("{")

	mv := metago.ValueOf(obj)
	var i int
	for _, mf := range mv.Fields() {
		if mf.Value().(time.Time).IsZero() {
			continue
		}

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := mf.Name()
		if v := mf.StructTagGet("json"); v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")

		switch v := mf.Value().(type) {
		case int64:
			buf.WriteString(strconv.FormatInt(v, 10))
		case int:
			buf.WriteString(strconv.Itoa(v))
		case string:
			buf.WriteString(strconv.Quote(v))
		case time.Time:
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}

		i++
	}

	buf.WriteString("}")

	s := buf.String()
	return *(*[]byte)(unsafe.Pointer(&s)), nil
}
