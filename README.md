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
数据集为某消费金融公司用户画像，包含是否是“坏”用户的标签，总数为5000行
### auc

### top-N accuracy

## TODO List
1. 针对iforest不支持category类型的改进
2. 实现cmd
3. 加入空值后，iforest稳定性变差，应考虑增加树深度
 
