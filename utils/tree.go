/*
@author '彼时思默'
@time 2020/3/25 14:06
@describe:
*/
package utils

import (
	"fmt"
	"strings"
)

type Tree struct {
	Name  string
	Nodes []*Tree
	Num   int
	deep  int
}

func NewTree(name string) *Tree {
	return &Tree{
		Name:  name,
		Nodes: make([]*Tree, 0),
		Num:   0,
	}
}

func (t *Tree) AddNodes(path string) {
	pathSp := t.path2array(path)
	n := t
	nn := t
	for i, name := range pathSp {
		nn = n.GetChildByName(name)
		if nn == nil {
			//fmt.Println("新建文件夹", Name)
			nodeTemp := NewTree(name)
			nodeTemp.deep = i + 1
			n.Nodes = append(n.Nodes, nodeTemp)
			n = n.GetChildByName(name)
		} else {
			n = nn
		}
	}

	t.RefreshNodeNum()
}

func (t Tree) RefreshNodeNum() {
	if t.HasChildren() {
		for _, t := range t.Nodes {
			t.RefreshNodeNum()
		}
	} else {
		t.Num++
		return
	}
	for _, c := range t.Nodes {
		t.Num += c.Num
	}
}

func (t Tree) HasChildren() bool {
	if len(t.Nodes) == 0 {
		return false
	}
	return true
}

func (t Tree) GetChildByName(name string) *Tree {
	if t.HasChildren() {
		for _, node := range t.Nodes {
			if node.Name == name {
				return node
			}
		}
	}
	return nil
}

//func (t Tree) SearchDepthFirst(Name string) *Tree {
//	chP := make(chan *Tree, 100)
//	chP <- &t
//	for target := range chP {
//		fmt.Println(target.deep, target.Name)
//		if target.Name == Name {
//			return target
//		}
//		if target.HasChildren() {
//			for _, target := range target.Nodes {
//				chP <- target
//			}
//		}
//		if len(chP) == 0 {
//			return nil
//		}
//	}
//	return nil
//}

func (t Tree) SearchBreadFirst(name string) *Tree {
	chP := make(chan *Tree, 100)
	chP <- &t
	for target := range chP {
		if target.Name == name {
			return target
		}
		if target.HasChildren() {
			for _, c := range target.Nodes {
				chP <- c
			}
		}
		if len(chP) == 0 {
			return nil
		}
	}
	return nil
}

func (t Tree) List(dirPath string) {
	dirSp := t.path2array(dirPath)
	n := &t
	for _, dir := range dirSp {
		if strings.Contains(dir, "/") {
			dir += "/"
		}
		dirName := n.Name
		n = n.GetChildByName(dir)
		if n == nil {
			fmt.Println("文件夹", dirName, "不存在")
			return
		}
	}
	for _, file := range n.Nodes {
		fmt.Println(file)
	}
}

func (t Tree) path2array(path string) []string {
	pathSp := strings.Split(path, "/")
	pathLen := len(pathSp)
	nodes := make([]string, 0)
	for i := 0; i < pathLen; i++ {
		name := pathSp[i]
		if pathLen >= 2 && i < pathLen-1 {
			nodes = append(nodes, name+"/")
		}
	}
	name := pathSp[pathLen-1]
	if len(name) > 0 {
		nodes = append(nodes, name)
	}
	return nodes
}

func (t Tree) PrintlnTree() {
	fmt.Println(t.Name)
	if t.HasChildren() {
		for _, t := range t.Nodes {
			t.PrintlnTree()
		}
	} else {
		return
	}
}
