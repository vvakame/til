// Code generated by metago. DO NOT EDIT.

//+build !metago

package benchmark

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/vvakame/til/go/metago"
)

// TODO 無くても動くように
var _ metago.Value = nil

type FooMetago struct {
	ID        int64
	Kind      string
	Name      string `json:"nickname"`
	Age       int
	CreatedAt time.Time
}

var propertyNameCache map[string]string

func (obj *FooMetago) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	buf.Grow(1024)
	if propertyNameCache == nil {
		propertyNameCache = make(map[string]string)
	}

	buf.WriteString("{")

	var i int
	{

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := "ID"
		if v := ""; v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")
		{

			buf.WriteString(strconv.FormatInt(obj.ID, 10))
		}

		i++
	}
	{

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := "Kind"
		if v := ""; v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")
		{

			buf.WriteString(strconv.Quote(obj.Kind))
		}

		i++
	}
	{

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := "Name"
		if v := "nickname"; v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")
		{

			buf.WriteString(strconv.Quote(obj.Name))
		}

		i++
	}
	{

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := "Age"
		if v := ""; v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")
		{

			buf.WriteString(strconv.Itoa(obj.Age))
		}

		i++
	}
	{
		if obj.CreatedAt.IsZero() {
			goto metagoGoto0

		}

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := "CreatedAt"
		if v := ""; v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")
		{

			b, err := obj.CreatedAt.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}

		i++
	}
metagoGoto0:
	;

	buf.WriteString("}")

	s := buf.String()
	return *(*[]byte)(unsafe.Pointer(&s)), nil
}