package main

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
	root      *Node
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
	model.SetParams(100, 256, F_NULL_MODE)
	return &model
}

func (f *IForest) SetParams(nTrees, subSample int, mode int) {
	f.nTrees = nTrees
	f.subSample = subSample
	f.mode = mode
}

func (f *IForest) Fit(rows []Row) {

}

func (f *IForest) Predict(rows []Row) (scores []floatt32) {

}

func (f *IForest) Evaluate(rows []Row) (auc float32) {

}
