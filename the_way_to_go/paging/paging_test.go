package paging_test

import (
	"github.com/grifree/code/the_way_to_go/paging"
	"log"
	"testing"
)

func TestCreateData(t *testing.T) {
	render := paging.CreateData(paging.Gen{
		Page:1,
		Total:-1,
	})
	log.Print(render)

}
