package paging_test

import (
	"github.com/grifree/code/the_way_to_go/paging"
	gis "github.com/og/x/test"
	"testing"
)

func TestGenCheckAndFix(t *testing.T) {
	is := gis.New(t)
	// gen.Page < 0
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.Page cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Page: -1,
		})
	}()
	// gen.Page == 0
	func() {
		// TODO 测试log打印内容
		_,gen := paging.CreateData(paging.Gen{
			Page: 0,
			PerPage:10,
			JumpPageInterval:1,
		})
		is.Eql(gen.Page, 1)
	}()
	// gen.Total < 0
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
	// gen.PerPage < 0
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
	// gen.PerPage == 0
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
	// gen.JumpPageInterval < 0
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
	// gen.JumpPageInterval == 0
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.JumpPageInterval cannot be 0")
		}()
		paging.CreateData(paging.Gen{
			Page:1,
			PerPage:1,
			JumpPageInterval:0,
		})
	}()
	// gen.ClosestPageLength < 0
	func() {
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "paging: paging.CreateData(gen) gen.ClosestPageLength cannot less zero")
		}()
		paging.CreateData(paging.Gen{
			Page:1,
			PerPage:1,
			ClosestPageLength:-1,
			JumpPageInterval:1,
		})
	}()
}

func TestRenderExistPaging(t *testing.T) {
	is := gis.New(t)
	// total = 0
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:0,
			PerPage:10,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ExistPaging, false)
	}()
	// total > 0
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:100,
			PerPage:10,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ExistPaging, true)
	}()
}

func TestRenderLastPage(t *testing.T){
	is := gis.New(t)
	// total = 0
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:0,
			PerPage:10,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.LastPage, 0)
	}()
	// 总页数刚好 整除无余
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:100,
			PerPage:10,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.LastPage, 10)
	}()
	// 有余
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:99,
			PerPage:10,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.LastPage, 10)
	}()
}

func TestRenderIsFirstOrLastPage(t *testing.T) {
	is := gis.New(t)
	// page > LastPage,  只1页
	func(){
		render,gen := paging.CreateData(paging.Gen{
			Page: 2,
			Total:1,
			PerPage:10,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(gen.Page, 1)
		is.Eql(render.LastPage, 1)
		is.Eql(render.IsFirstPage, true)
		is.Eql(render.IsLastPage, true)
	}()
	// 首页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:2,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.LastPage, 2)
		is.Eql(render.IsFirstPage, true)
		is.Eql(render.IsLastPage, false)
	}()
	// 最后一页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:2,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.IsFirstPage, false)
		is.Eql(render.IsLastPage, true)
	}()
	// 中间页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:3,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.IsFirstPage, false)
		is.Eql(render.IsLastPage, false)
	}()
}

func TestRenderClosestPagePrev(t *testing.T) {
	is := gis.New(t)
	// page = 1 没有向前页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:1,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
	}()
	// 有向前页1个, 但配置ClosestPageLength=0
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:2,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
	}()
	// 有向前页多个, 但配置ClosestPageLength=0
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 3,
			Total:3,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
	}()
	// 首页与当前页间隔 < 配置ClosestPageLength
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 3,
			Total:3,
			PerPage:1,
			ClosestPageLength:2,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{2})
	}()
	// 首页与当前页间隔 == 配置ClosestPageLength
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 4,
			Total:4,
			PerPage:1,
			ClosestPageLength:2,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{2,3})
	}()
	// 首页与当前页间隔 > 配置ClosestPageLength
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:2,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{3,4})
	}()
}

func TestRenderClosestPageNext(t *testing.T) {
	is := gis.New(t)
	// 总共1页, 没有向后页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:1,
			PerPage:1,
			ClosestPageLength:1,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	// 有向后页, 但page与lastpage间隔不足
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:2,
			PerPage:1,
			ClosestPageLength:1,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	// page与lastpage间隔足够, 但配置ClosestPageLength=0
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:3,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	// page与lastpage间隔 < 配置ClosestPageLength
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:3,
			PerPage:1,
			ClosestPageLength:2,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{2})
	}()
	// page与lastpage间隔 == 配置ClosestPageLength
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:4,
			PerPage:1,
			ClosestPageLength:2,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{2,3})
	}()
	// page与lastpage间隔 > 配置ClosestPageLength
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:5,
			PerPage:1,
			ClosestPageLength:2,
			JumpPageInterval:1,
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
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, true)
		is.Eql(render.ClosestPage.Next.PageList, []int{4,5,6})
	}()
	// 最后1页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 9,
			Total:9,
			PerPage:1,
			ClosestPageLength:3,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
	// 倒数第2页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 8,
			Total:9,
			PerPage:1,
			ClosestPageLength:3,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Next.Exist, false)
		is.Eql(render.ClosestPage.Next.PageList, nil)
	}()
}

func TestRenderJumpPagePrev(t *testing.T) {
	is := gis.New(t)
	// 总1页
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 1,
			Total:1,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		// TODO 怎么判断结构体相等
		is.Eql(render.JumpPage.Prev.Exist, false)
		is.Eql(render.JumpPage.Prev.Interval, 0)
	}()
	// firstPage与page间隔不足
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 2,
			Total:5,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
		is.Eql(render.JumpPage.Prev.Exist, false)
		is.Eql(render.JumpPage.Prev.Interval, 0)
	}()
	// firstPage与page间隔, 去除ClosestPage后, 空间不足
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:3,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{2,3,4})
		is.Eql(render.JumpPage.Prev.Exist, false)
		is.Eql(render.JumpPage.Prev.Interval, 0)
	}()
	// firstPage与page间隔足
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 3,
			Total:5,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
		is.Eql(render.JumpPage.Prev.Exist, true)
		is.Eql(render.JumpPage.Prev.Interval, 2)
	}()
	// firstPage与page间隔足, 无ClosestPage
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:0,
			JumpPageInterval:1,
		})
		is.Eql(render.ClosestPage.Prev.Exist, false)
		is.Eql(render.ClosestPage.Prev.PageList, nil)
		is.Eql(render.JumpPage.Prev.Exist, true)
		is.Eql(render.JumpPage.Prev.Interval, 4)
	}()
	// firstPage与page间隔足, 有ClosestPage
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:1,
			JumpPageInterval:2,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{4})
		is.Eql(render.JumpPage.Prev.Exist, true)
		is.Eql(render.JumpPage.Prev.Interval, 3)
	}()
	// JumpPageInterval 刚好至最远 可跳空间
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:1,
			JumpPageInterval:5,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{4})
		is.Eql(render.JumpPage.Prev.Exist, true)
		is.Eql(render.JumpPage.Prev.Interval, 1)
	}()
	// JumpPageInterval 超出 可跳空间
	func(){
		render,_ := paging.CreateData(paging.Gen{
			Page: 5,
			Total:5,
			PerPage:1,
			ClosestPageLength:1,
			JumpPageInterval:10,
		})
		is.Eql(render.ClosestPage.Prev.Exist, true)
		is.Eql(render.ClosestPage.Prev.PageList, []int{4})
		is.Eql(render.JumpPage.Prev.Exist, true)
		is.Eql(render.JumpPage.Prev.Interval, 1)
	}()
}

func TestRenderJumpPageNext(t *testing.T) {
	is := gis.New(t)
	// 只有1页
	func() {
		render, _ := paging.CreateData(paging.Gen{
			Page:              1,
			Total:             1,
			PerPage:           1,
			ClosestPageLength: 0,
			JumpPageInterval:  1,
		})
		is.Eql(render.JumpPage.Next.Exist, false)
		is.Eql(render.JumpPage.Next.Interval, 0)
	}()
	// 最后1页
	func() {
		render, _ := paging.CreateData(paging.Gen{
			Page:              5,
			Total:             5,
			PerPage:           1,
			ClosestPageLength: 0,
			JumpPageInterval:  1,
		})
		is.Eql(render.JumpPage.Next.Exist, false)
		is.Eql(render.JumpPage.Next.Interval, 0)
	}()
	// page与lastPage间隔不足
	func() {
		render, _ := paging.CreateData(paging.Gen{
			Page:              4,
			Total:             5,
			PerPage:           1,
			ClosestPageLength: 0,
			JumpPageInterval:  1,
		})
		is.Eql(render.JumpPage.Next.Exist, false)
		is.Eql(render.JumpPage.Next.Interval, 0)
	}()
	// page与lastPage间隔1
	func() {
		render, _ := paging.CreateData(paging.Gen{
			Page:              3,
			Total:             5,
			PerPage:           1,
			ClosestPageLength: 0,
			JumpPageInterval:  1,
		})
		is.Eql(render.JumpPage.Next.Exist, true)
		is.Eql(render.JumpPage.Next.Interval, 4)
	}()
	// JumpInterval 大于 page与lastPage间隔
	func() {
		render, _ := paging.CreateData(paging.Gen{
			Page:              1,
			Total:             5,
			PerPage:           1,
			ClosestPageLength: 0,
			JumpPageInterval:  10,
		})
		is.Eql(render.JumpPage.Next.Exist, true)
		is.Eql(render.JumpPage.Next.Interval, 5)
	}()
}