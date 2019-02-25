package main

import (
	"log"
	"testing"
)

func TestLoadData(t *testing.T) {
	//path := "./data/test.csv"
	//
	//labels, rows := LoadData(path, false, 0, true, 0, 1)
	//log.Println(labels)
	//log.Println(rows)
	//
	//labels0, rows0 := LoadData(path, true, 1., true, 0, 1)
	//log.Println(labels0)
	//log.Println(rows0)
	//log.Println(len(rows0))

	path_large := "./data/5000.csv"

	labels1, rows1 := LoadData(path_large, false, 0, true, 0, 1)
	//log.Println(labels1)
	//log.Println(rows1)
	log.Println(len(labels1))
	log.Println(len(rows1), len(rows1[0]))
}
