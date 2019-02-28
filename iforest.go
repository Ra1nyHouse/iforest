package main

import (
	"log"
	"math"
	"math/rand"
)

const (
	// 	允许特征包含空值
	F_NULL_MODE = 0x1

	// 允许特征包含类别类型
	F_CAT_MODE = 0x2
)

type IForest struct {
	nTrees    int
	subSample int
	mode      int
	randSeed  int
	roots     []*Node
}

type Node struct {
	Left       *Node
	Right      *Node
	SplitAtt   int
	SplitValue Elem
	Size       int
}

func NewModel() *IForest {
	model := IForest{}
	model.SetParams(100, 256, F_NULL_MODE, 1993)
	return &model
}

func (f *IForest) SetParams(nTrees, subSample int, mode int, randSeed int) {
	f.nTrees = nTrees
	f.subSample = subSample
	f.mode = mode
	f.randSeed = randSeed
}

func (f *IForest) Fit(rows []Row) {
	r := rand.New(rand.NewSource(int64(f.randSeed)))
	l := int(math.Ceil(math.Log2(float64(f.subSample))))
	ch := make(chan *Node, 100)
	for i := 0; i < f.nTrees; i++ {
		// 有放回采样
		row_index := make([]int, f.subSample)
		for j := 0; j < f.subSample; j++ {
			row_index[j] = r.Intn(len(rows))
		}
		go func(treeID int, rSeed int64) {
			ch <- _itree(rows, row_index, 0, l, rand.New(rand.NewSource(rSeed)))
			log.Printf("goroutine: iTree #%d done...\n", treeID)
		}(i, r.Int63n(10000))
	}

	roots := make([]*Node, 0, 1000)
	for i := 0; i < f.nTrees; i++ {
		roots = append(roots, <-ch)
		//log.Printf("main: iTree #%d done...", i)
	}
	f.roots = roots

	log.Println("main: finish fit...")
}

func _itree(rows []Row, row_index []int, e int, l int, r *rand.Rand) *Node {
	if e >= l || len(row_index) <= 1 {
		return &Node{Size: len(row_index)}
	} else {
		max_k := 5
		for k := 0; k < max_k; k++ {
			q := r.Intn(len(rows[0]))
			not_null_index := make([]int, 0, 1000)
			min_v, max_v := float32(0), float32(0)
			has_assign := false

			// 遍历获取最大最小值，同时获取非空值索引
			for i := 0; i < len(row_index); i++ {
				if elem := rows[row_index[i]][q]; elem.Valid {
					not_null_index = append(not_null_index, row_index[i])
					if !has_assign {
						has_assign = true
						min_v = elem.Value
						max_v = elem.Value
					} else {
						if min_v > elem.Value {
							min_v = elem.Value
						}
						if max_v < elem.Value {
							max_v = elem.Value
						}
					}
				}
			}

			// 当前选择的q列全为空值，没有足够的信息进行划分
			// 注意至多进行max_k次重划分, 超过上限则停止划分
			if len(not_null_index) == 0 {
				continue
			}

			p := float32(min_v + (max_v-min_v)*r.Float32())
			left_row_index := make([]int, 0, 1000)
			right_row_index := make([]int, 0, 1000)
			for i := 0; i < len(row_index); i++ {
				// 如果值为空值，就从非空值中随机选一个作为替代
				if elem := rows[row_index[i]][q]; !elem.Valid {
					if rand_elem := rows[not_null_index[r.Intn(len(not_null_index))]][q]; rand_elem.Value < p {
						left_row_index = append(left_row_index, row_index[i])
					} else if rand_elem.Value >= p {
						right_row_index = append(right_row_index, row_index[i])
					}
				} else if elem.Value < p {
					left_row_index = append(left_row_index, row_index[i])
				} else if elem.Value >= p {
					right_row_index = append(right_row_index, row_index[i])
				}
			}

			return &Node{Left: _itree(rows, left_row_index, e+1, l, r),
				Right:      _itree(rows, right_row_index, e+1, l, r),
				SplitAtt:   q,
				SplitValue: Elem{Value: p, Valid: true},
				Size:       len(row_index)}
		}
	}
	// 至多进行max_k次重划分, 超过上限则停止划分
	return &Node{Size: len(row_index)}
}

func (f *IForest) Predict(rows []Row) []float32 {
	ch := make(chan float32, 100)
	scores := make([]float32, 0, 1000)
	r := rand.New(rand.NewSource(int64(f.randSeed)))
	for _, row := range rows {
		for _, node := range f.roots {
			// 在树的粒度上并行
			//go func(r *rand.Rand) {
			//	ch <- _path(row, node, 0, r)
			//}(rand.New(rand.NewSource(r.Int63n(10000))))
			//go func() {
			//	ch <- _path(row, node, 0, rand.New(rand.NewSource(r.Int63n(10000))))
			//}()
			go func(rSeed int64) {
				ch <- _path(row, node, 0, rand.New(rand.NewSource(rSeed)))
			}(r.Int63n(10000))
		}
		score := float32(0)
		for i := 0; i < len(f.roots); i++ {
			score += <-ch
		}
		score = score / float32(len(f.roots))
		scores = append(scores, score)
	}
	return scores
}

// 对应论文中的PathLength函数
func _path(row Row, node *Node, e int, r *rand.Rand) float32 {
	if node.Left == nil || node.Right == nil {
		return float32(e) + _c(node.Size)
	}
	elem := row[node.SplitAtt]
	// 正常情况，值没有缺失
	if elem.Valid {
		if elem.Value < node.SplitValue.Value {
			return _path(row, node.Left, e+1, r)
		} else {
			return _path(row, node.Right, e+1, r)
		}
	} else {
		// 值缺失，根据树的左右节点数量，随机到左孩子或右孩子节点
		leftN := node.Left.Size
		rightN := node.Right.Size
		randN := r.Intn(leftN + rightN)
		if randN < leftN {
			return _path(row, node.Left, e+1, r)
		} else {
			return _path(row, node.Right, e+1, r)
		}
	}
}

// 对应论文中的c函数
func _c(size int) float32 {
	if size < 2 {
		return 0.
	} else {
		return 2.*(float32(math.Log(float64(size-1)))+0.5772156649) - 2.*(float32(size)-1.)/float32(size)
	}
}

func (f *IForest) Evaluate(rows []Row, labels []int8) float32 {
	scores := f.Predict(rows)
	// label 1表示坏用户， 而score越小表示用户越“坏”
	for i, score := range scores {
		scores[i] = -score
	}
	auc := float32(0)
	mPos := 0
	mNeg := 0
	for i, labeli := range labels {
		if labeli == 0 {
			mNeg += 1
			continue
		}
		mPos += 1
		for j, labelj := range labels {
			if labelj == 1 {
				continue
			}
			if scores[i] < scores[j] {
				auc += 1
			} else if scores[i] == scores[j] {
				auc += 0.5
			}
		}
	}
	auc = auc / float32(mPos*mNeg)
	auc = 1 - auc
	return auc
}
