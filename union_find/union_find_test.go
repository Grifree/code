package union_find_test

import (
	"log"
	gconv "github.com/og/x/conv"
	"testing"
)

type QuickFind struct {
	id []int // []int 中 index 表示点 value 表示分量id 可以理解为 []int 类似 map[pointer]componentID
	count int // 分量数量

}
func NewQuickFind(n int) QuickFind {
	uf := QuickFind{}
	uf.id = []int{}
	uf.count = n
	for i:=0;i<n;i++ {
		uf.id = append(uf.id, i)
	}
	return uf
}
func (uf QuickFind) Connected(p int, q int) bool {
	return uf.Find(p) == uf.Find(q)
}
func (uf QuickFind) Find(p int) int {
	return uf.id[p]
}
func (uf *QuickFind) Union(p int, q int) {
	pID := uf.Find(p)
	qID := uf.Find(q)
	if pID == qID {return}
	for i:=0;i< len(uf.id) ; i++ {
		if uf.id[i] == pID {
			uf.id[i] = qID
			uf.count--
		}
	}
}

func (uf QuickFind) Count() int {
	return uf.count
}
func (uf QuickFind) String() string {
	data := "\r\n"
	groupMap := map[int]string{}
	for value, componentID := range uf.id {
		groupMap[componentID] += gconv.IntString(value) + " "
	}
	for componentID, message := range groupMap {
		data += gconv.IntString(componentID) + " [ " + message + "]\r\n"
	}
	return data
}

func TestNewQuickFind(t *testing.T) {
	a := NewQuickFind(9)
	log.Print(a.String())
	a.Union(5,6)
	log.Print(a.String())
	a.Union(6,8)
	log.Print(a.String())
}