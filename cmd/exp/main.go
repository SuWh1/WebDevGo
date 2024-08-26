package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
	Bio  string
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "dauren",
		Age:  18,
		Meta: UserMeta{
			Visits: 4,
		},
		Bio: `<script>alert("haha you v been joked");</script>`,
	}

	err = t.Execute(os.Stdout, user) // execute - take a template and proccess it
	if err != nil {
		panic(err)
	}
}
