package main

import (
	"fmt"
	"reflect"
	"encoding/json"
	"errors"
	"strconv"
)
const (
	TypeA = "a identifier"
	TypeB = "b identifier"
)

type A struct {
	Type string
	A    int64
}

func (a *A) GetKey() string {
	return strconv.FormatInt(a.A, 10)
}

func (b *B) GetKey() string {
	return b.S
}

type X interface {
	GetKey() string
}
type B struct {
	Type string
	S    string
}

func newA() *A {
	return &A{
		Type: TypeA,
	}
}

func newB() *B {
	return &B{
		Type: TypeB,
	}
}

type JSONMap = map[string]interface{}

func Save(x []X) []byte {
	bytes, _ := json.Marshal(x)
	return bytes
}

func Retrieve(bytes []byte) []X {
	var p []JSONMap
	err := json.Unmarshal(bytes, &p)
	if err != nil {
		fmt.Println(err)
	}

	res := make([]X, 0)
	for _, m := range p {
		if t, ok := m["Type"]; ok {
			if str, isStr := t.(string); isStr {
				switch str {
				case TypeA:
					var a A
					for k, v := range m {
						err := SetField(&a, k, v)
						if err != nil {
							fmt.Println("Error", err)
						}
					}
				case TypeB:
					var b B
					for k, v := range m {
						err := SetField(&b, k, v)
						if err != nil {
							fmt.Println("Error", err)
						}
					}
				default:
					fmt.Println("Unknown type", str)
					continue
				}
			}
			// here we convert the map to the struct

		}
		fmt.Println("Unknown struct", m)
	}
	return res
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func main() {

	a := newA()
	a.A = 86

	b := newB()
	b.S = "8"

	x := []X{
		a,
		b,
	}

	bytes := Save(x)

	t := Retrieve(bytes)

	fmt.Println(len(t))
}