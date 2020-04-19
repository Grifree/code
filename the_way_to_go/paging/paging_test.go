package paging_test

import (
	"github.com/grifree/code/the_way_to_go/paging"
	gis "github.com/og/x/test"
	"testing"
)

func TestGenCheckAndFix(t *testing.T) {
	is := gis.New(t)
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.Page cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Page: -1,
		})
	}()
	func() {
		// TODO 测试log打印内容
		_,gen := paging.CreateData(paging.Gen{
			Page: 0,
			PerPage:10,
		})
		is.Eql(gen.Page, 1)
	}()
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.Total cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Total: -2,
			Page:1,
		})
	}()
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.PerPage cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Total: 1,
			Page:1,
			PerPage:-1,
		})
	}()
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.PerPage cannot be 0")
		}()
		paging.CreateData(paging.Gen{
			Total: 1,
			Page:1,
			PerPage:0,
		})
	}()
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.JumpPageInterval cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Page:1,
			PerPage:1,
			JumpPageInterval:-1,
		})
	}()
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.ClosestPageLength cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Page:1,
			PerPage:1,
			ClosestPageLength:-1,
		})
	}()
}

func TestRenderExistPaging(t *testing.T) {
	is := gis.New(t)
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:0,
			PerPage:10,
		})
		is.Eql(render.ExistPaging, false)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:100,
			PerPage:10,
		})
		is.Eql(render.ExistPaging, true)
	}()
}

func TestRenderLastPage(t *testing.T){
	is := gis.New(t)
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:0,
			PerPage:10,
		})
		is.Eql(render.LastPage, 0)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:100,
			PerPage:10,
		})
		is.Eql(render.LastPage, 10)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:99,
			PerPage:10,
		})
		is.Eql(render.LastPage, 10)
	}()
}

func TestRenderIsFirstOrLastPage(t *testing.T) {
	is := gis.New(t)
	func(){
		render,gen := paging.CreateData(paging.Gen{
			Page: 2,
			Total:1,
			PerPage:10,
		})
		is.Eql(gen.Page, 1)
		is.Eql(render.IsFirstPage, true)
		is.Eql(render.IsLastPage, true)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:2,
			PerPage:1,
		})
		is.Eql(render.IsFirstPage, true)
		is.Eql(render.IsLastPage, false)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:2,
			PerPage:1,
		})
		is.Eql(render.IsFirstPage, false)
		is.Eql(render.IsLastPage, true)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:3,
			PerPage:1,
		})
		is.Eql(render.IsFirstPage, false)
		is.Eql(render.IsLastPage, false)
	}()
}

func TestRenderClosestPagePrev(t *testing.T) {
	is := gis.New(t)
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:1,
			PerPage:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:2,
			PerPage:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 3,
			Total:3,
			PerPage:1,
			ClosestPageLength:0,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 3,
			Total:3,
			PerPage:1,
			ClosestPageLength:2,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{2})
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 4,
			Total:4,
			PerPage:1,
			ClosestPageLength:2,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{2,3})
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:2,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{3,4})
	}()
}


func TestRenderClosestPageNext(t *testing.T) {
	is := gis.New(t)
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:1,
			PerPage:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:2,
			PerPage:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:3,
			PerPage:1,
			ClosestPageLength:0,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:3,
			PerPage:1,
			ClosestPageLength:2,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{2})
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:4,
			PerPage:1,
			ClosestPageLength:2,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{2,3})
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:5,
			PerPage:1,
			ClosestPageLength:2,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{2,3})
	}()
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 3,
			Total:9,
			PerPage:1,
			ClosestPageLength:3,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{4,5,6})
	}()
}