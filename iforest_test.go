package main

import (
	"fmt"
	"testing"
)

func ExampleAb() {
	a := [8]int{}
	fmt.Println(a)
}

func Test_Division_1(t *testing.T) {
	a := make([][3]int, 3, 5) // len(b)=0, cap(b)=5
	fmt.Println(a)
	fmt.Println(a[0])
	fmt.Println(a[0][1])

	b := [3]int{0, 1, 2}
	fmt.Println(b)

	c := make([]interface{}, 5)
	c[1] = 20
	fmt.Print(c)
}
