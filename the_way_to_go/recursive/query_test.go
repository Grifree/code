package recursive_test

import (
	"github.com/grifree/code/the_way_to_go/recursive"
	gtest "github.com/og/x/test"
	"testing"
)

func TestFindAllChild(t *testing.T) {
	as := gtest.NewAS(t)
	as.Equal(
		recursive.FindAllChild("1"),
		[]string{"1","1-1","1-1-1","1-1-2","1-1-1-1","1-1-1-2"},
	)
}
