package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id   int
	Name string
	age  int
}

type Manager struct {
	User
	title string
}

func (u User) Hello() {
	fmt.Println("hello world!")
}

func main() {
	u := Manager{User: User{11, "ok", 20}, title: "title222"}
	Info(u)
}

func Info(o interface{}) {
	t := reflect.TypeOf(o)
	fmt.Println(t.Name(), t.NumField(), t.Field(0))

	m := t.Method(0)
	fmt.Println(m.Name)
	fmt.Println(m.Type)
	fmt.Println(m.Func)
	fmt.Println(m.Index)
	fmt.Println(m.PkgPath)

	fmt.Println("kind", t.Kind())

	v := reflect.ValueOf(o)
	fmt.Println(v.NumField(), v.Field(0) )

	/*v:=reflect.ValueOf(o)
	fmt.Println("fileds")

	for i:=0;i<t.NumField();i++{
		fmt.Println(t.Field(i))
		fmt.Println(v.Field(i).Interface())
	}*/
}
