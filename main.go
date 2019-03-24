package main

import (
	"log"
	"math/rand"
)

func main() {
	path := "./data/5000.csv"
	time := 50
	labels, rows := LoadData(path, false, 0, true, 0, 1)
	//labels, rows := LoadData(path, true, -1, true, 0, 1)

	aucs := make([]float32, time)
	top20s := make([]float32, time)
	top50s := make([]float32, time)
	top100s := make([]float32, time)
	top200s := make([]float32, time)
	top2000s := make([]float32, time)

	for i := 0; i < time; i++ {
		f := NewModel()
		//f.SetParams(1000, 512, F_NULL_MODE, 1993)
		f.SetParams(100, 256, F_NULL_MODE, rand.Int())
		f.Fit(rows)

		scores := f.Predict(rows)
		aucs[i] = f.Metric(scores, labels, AUC)
		top20s[i] = f.Metric(scores, labels, Top20Pre)
		top50s[i] = f.Metric(scores, labels, Top50Pre)
		top100s[i] = f.Metric(scores, labels, Top100Pre)
		top200s[i] = f.Metric(scores, labels, Top200Pre)
		top2000s[i] = f.Metric(scores, labels, Top2000Pre)

	}

	log.Println("AUC:", avg(aucs))
	log.Println("Top20 Pre.:", avg(top20s))
	log.Println("Top50 Pre.:", avg(top50s))
	log.Println("Top100 Pre.:", avg(top100s))
	log.Println("Top200 Pre.:", avg(top200s))
	log.Println("Top2000 Pre.:", avg(top2000s))

}

func avg(metrics []float32) float32 {
	s := float32(0)
	for _, m := range metrics {
		s += m
	}
	return s / float32(len(metrics))
}

func init() {
	// LogFlag 日志显示格式
	lf := log.Ldate | log.Ltime | log.Lshortfile
	log.SetFlags(lf)
}
