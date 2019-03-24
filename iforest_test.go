package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
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
	labels, rows := LoadData(path, false, 0, true, 0, 1)

	f := NewModel()
	f.Fit(rows)

	seeds := make([]int, 100)
	for i, node := range f.roots {
		//log.Println(node.TreeSeed)
		seeds[i] = node.TreeSeed

		if node.TreeSeed == 79 {
			log.Println(i)
			f.Show(i)

			p := _path(rows[0], node, 0, rand.New(rand.NewSource(int64(500))))
			log.Println(p)
		}
	}
	sort.Ints(seeds)
	log.Println(seeds)

	fmt.Println(f.Predict(rows))
	fmt.Println(f.Evaluate(rows, labels, AUC))

	fmt.Println("--")
	a := float32(0)
	for _, node := range f.roots {
		b := _path(rows[0], node, 0, rand.New(rand.NewSource(int64(node.TreeSeed))))
		fmt.Println(b)
		a += b
	}
	fmt.Println(a)

	fmt.Println("----")
	ch := make(chan float32, 100)
	for _, node := range f.roots {
		// 在树的粒度上并行
		//go func(rSeed int64) {
		//	ch <- _path(row, node, 0, rand.New(rand.NewSource(rSeed)))
		//}(r.Int63n(10000))
		go func(rSeed int, n *Node) {
			p := _path(rows[0], n, 0, rand.New(rand.NewSource(int64(rSeed))))
			fmt.Println("go:", p, " ,node seed:", node.TreeSeed)
			ch <- p
		}(node.TreeSeed, node)
	}
	score := float32(0)
	for i := 0; i < len(f.roots); i++ {
		b := <-ch
		fmt.Println(b)
		score += b
	}
	fmt.Println("score:", score)

}

func TestIForest_Fit2(t *testing.T) {
	path := "./data/5000.csv"
	_, rows := LoadData(path, false, 0, true, 0, 1)

	f := NewModel()
	f.Fit(rows)

	ch := make(chan float32, 100)
	for _, node := range f.roots {
		// 在树的粒度上并行
		//go func(rSeed int64) {
		//	ch <- _path(row, node, 0, rand.New(rand.NewSource(rSeed)))
		//}(r.Int63n(10000))
		go func(rSeed int, n *Node) {
			p := _path(rows[0], n, 0, rand.New(rand.NewSource(int64(rSeed))))
			fmt.Println("go:", p, " ,node seed:", rSeed)
			ch <- p
		}(node.TreeSeed, node)
	}
	score := float32(0)
	for i := 0; i < len(f.roots); i++ {
		b := <-ch
		fmt.Println(b)
		score += b
	}
	fmt.Println("score:", score)
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

	//fmt.Println(f.Evaluate(rows, labels, AUC))
	fmt.Println(f.Evaluate(rows, labels, Top100Pre))
}

func TestIForest_Evaluate2(t *testing.T) {
	path := "./data/5000.csv"
	labels, rows := LoadData(path, false, 0, true, 0, 1)
	//labels, rows := LoadData(path, true, -1, true, 0, 1)
	f := NewModel()
	f.SetParams(100, 256, F_NULL_MODE, 1993)
	f.Fit(rows)

	fmt.Println(f.Evaluate(rows, labels, AUC))
	fmt.Println(f.Evaluate(rows, labels, AUC))

	scores := f.Predict(rows)
	fmt.Println(f.Metric(scores, labels, AUC))
	//fmt.Println(f.Evaluate(rows, labels, Top100Pre))
}

func TestSort(t *testing.T) {
	scores := []float32{23, 0.1, 0.5, 1.1, 100}
	indices := make([]struct {
		I int
		V float32
	}, len(scores))
	for i, v := range scores {
		indices[i].I = i
		indices[i].V = v
	}
	sort.Slice(indices, func(i, j int) bool {
		return indices[i].V < indices[j].V
	})
	fmt.Println(indices)

	//index := make([]int, len(scores))
	//for i:=0 ; i<len(index) ; i++ {
	//	index[i] = i
	//}
	//sort.Slice(index, func(i, j inbool {
	//	return scores[i] < scores[j]
	//})
	//
	//fmt.Println(scores)
	//fmt.Println(index)
}

func TestRand(t *testing.T) {
	fmt.Println(rand.Intn(10000))
	fmt.Println(rand.Intn(10000))
	fmt.Println(rand.Intn(10000))
}

func TestAuc(t *testing.T) {
	scores := []float32{0.1, 0.4, 0.35, 0.8}
	labels := []int8{0, 0, 1, 1}

	// 输出0.75
	fmt.Println(auc(scores, labels))
}
