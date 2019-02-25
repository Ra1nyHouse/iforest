package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

//type Dataset struct {
//}
//
//type DataIter struct {
//
//}
//
//type Elementer interface {
//
//}

type Elem struct {
	Value float32
	Valid bool
}

type Row []Elem

// LoadData 载入数据
func LoadData(path string, replaceNull bool, nullValue float32, ignoreFirst bool, idCol int, labelCol int) ([]int8, []Row) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	if ignoreFirst {
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	rows := make([]Row, 0, 1000)
	labels := make([]int8, 0, 1000)
	for scanner.Scan() {
		line := scanner.Text()
		//log.Println(line)
		line = strings.Trim(line, "\n")

		row := make([]Elem, 0, 1000)
		for i, chs := range strings.Split(line, ",") {
			if i == idCol {
				continue
			} else if i == labelCol {
				l, err := strconv.Atoi(chs)
				if err != nil {
					log.Fatal(err)
				}
				labels = append(labels, int8(l))
			} else {
				if chs == "" {
					if replaceNull {
						row = append(row, Elem{nullValue, true})
					} else {
						row = append(row, Elem{Valid: false})
					}
				} else {
					v, err := strconv.ParseFloat(chs, 32)
					if err != nil {
						log.Fatal(err)
					}
					row = append(row, Elem{float32(v), true})
				}
			}
		}
		//log.Println(row)
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
	return labels, rows
}
