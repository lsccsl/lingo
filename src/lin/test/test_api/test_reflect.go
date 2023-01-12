package main

import (
	"fmt"
	"reflect"
)

func F(callback interface{}, args ...interface{}) {
	v := reflect.ValueOf(callback)
	if v.Kind() != reflect.Func {
		panic("not a function")
	}
	vargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vargs[i] = reflect.ValueOf(arg)
	}

	vrets := v.Call(vargs)

	fmt.Print("\tReturn values: ", vrets)
}

func CB1() {
	fmt.Println("CB1 called")
}

func CB2() bool {
	fmt.Println("CB2 called")
	return false
}

func CB3(s string) {
	fmt.Println("CB3 called")
}

func CB4(s string) (bool, int) {
	fmt.Println("CB4 called")
	return false, 1
}

type TestR struct {
	a int `test:"aa",test2:"bb"`
}

func test_reflect() {

	ts := TestR{}

	t := reflect.TypeOf(ts)
	f, b := t.FieldByName("a")
	fmt.Println(f, b)

	F(CB1)
	F(CB2)
	F(CB3, "xxx")
	F(CB4, "yyy")
}
