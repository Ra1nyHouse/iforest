# iforest go实现

## iforest
1. 本项目基于go语言实现了 iforest （Isolation Forest， Fei Tony Liu, Kai Ming Ting， https://cs.nju.edu.cn/zhouzh/zhouzh.files/publication/icdm08b.pdf）
2. 使用了go routine，在**树的粒度上并行**（即并行的构建树）
3. 论文原文实现的iforest有两个缺陷：（I）不支持空值，（II）不支持category类型
4. 针对缺陷（I），本项目进行了改进

## 针对iforest不支持空值的改进
### 树的构建阶段（训练）
（数学符号表示沿用了论文）
1. 随机选择q列，如果样本x的q列为空，这个样本也会划分到左节点或者右节点，策略为：从当前数据集X随机选取一个q列不是空值得样本，将其q列值作为x的q列值得替代，然后和p值比较，决定划分到左节点或者右节点
2. 异常情况：如果当前样本q列值都为空，就重新选择q列，重复5次，q列值仍然为空，就终止构造树
### 计算样本的PathLength（测试）
1. 当前的划分列为splitAtt，如果样本x的splitAtt列值为空，样本仍然需要被划分到左节点或右节点，策略为：假设左节点所有的孩子节点|X|值为a，右节点所有的孩子节点|X|值为b，产生一个在[0,a+b）之间的随机数c，如果c<a，就将x划分到左节点，否则划分到右节点

## 实验结果
### 数据集
数据集为某消费金融公司用户画像，包含是否是“坏”用户的标签，总数为5000行，维度是1400维
### 结果

| 指标 | AUC | top50 Pre. | top100 Pre. | top200 Pre.|
| ------ | ------ | ------ | ------ | ------ |
| Standard Mode | 0.474 | 0.232 | 0.276 | 0.301 |
| NULL Mode | 0.536 | 0.331 | 0.362 | 0.389 |

Standard Mode与论文实现相同，遇到空值则填充-1；
NULL Mode针对空值做了改进。
所有实验重复20次取均值。
从结果来说支持空值的（NULL Mode）iforest表现更好，但不严谨。
一方面数据集特征处理比较简单，且特征不一定与金融属性呈强相关性；
一方面还需要理论论证解决空值策略的合理性。


## BUG

### 子树顺序不一致问题
树并行训练，因此构建每棵树结束时间无法确定，每次运行Fit产生的iforest都不同（子树顺序不同）。
不同顺序的子树会导致预测阶段结果的不同，因此需要为每颗子树绑定一个随机数生成对象。

### go程避免引用局部变量
错误示例如下：
```go
for _, node := range f.roots {
    // 在树的粒度上并行
    go func(_node *Node) {
        ch <- _path(row, _node, 0, rand.New(rand.NewSource(int64(_node.TreeSeed))))
    }(node)

    // 一定要注意协程不要引用局部变量，下面会导致BUG，node几乎是同一个
    //go func(rSeed int) {
    //	ch <- _path(row, node, 0, rand.New(rand.NewSource(int64(rSeed))))
    //}(node.TreeSeed)
}
```


## Trick
### 提升测试阶段效率
树的构建阶段（训练阶段）花费时间较少，主要时间消耗在测试阶段。
测试阶段每读入一个样例，启动多个协程分别计算样例在不同树上的path，
完成后关闭协程，继续读入样例，这样效率不高。

为了减少开启关闭协程的操作，可以不关闭每棵树对应的协程，而是阻塞等待后面的数据进入channel，并用一个协程持续读入数据。


## TODO List
1. 针对iforest不支持category类型的改进
2. 实现cmd接口
3. 加入空值后，iforest稳定性变差，应考虑增加树深度
 
