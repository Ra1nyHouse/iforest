package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleAb() {
	a := [8]int{}
	fmt.Println(a)
}

func Test_Division_1(t *testing.T) {
	for _, value := range rand.Perm(8) {
		fmt.Println(value)
	}

}

func TestIForest_Fit(t *testing.T) {
	path := "./data/5000.csv"
	_, rows := LoadData(path, false, 0, true, 0, 1)

	f := NewModel()
	f.Fit(rows)
}

func Test_C(t *testing.T) {
	fmt.Println(_c(1))
	fmt.Println(_c(2))
	fmt.Println(_c(3))
	fmt.Println(_c(4))
}

func TestIForest_Predict(t *testing.T) {
	path := "./data/5000.csv"
	_, rows := LoadData(path, false, 0, true, 0, 1)

	f := NewModel()
	//f.SetParams(1000, 512, F_NULL_MODE, 1993)
	f.Fit(rows)

	scores := f.Predict(rows)
	fmt.Println(len(scores))
	fmt.Println(scores[0:10])

}

func TestIForest_Evaluate(t *testing.T) {
	path := "./data/5000.csv"
	labels, rows := LoadData(path, false, 0, true, 0, 1)
	//labels, rows := LoadData(path, true, -1, true, 0, 1)
	f := NewModel()
	f.SetParams(1000, 512, F_NULL_MODE, 1993)
	f.Fit(rows)

	fmt.Println(f.Evaluate(rows, labels))
}
